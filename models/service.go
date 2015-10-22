package models

import (
	"errors"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	. "github.com/kr/pretty"
	"github.com/telamon/wharfmaster/util"
	"os"
	"regexp"
	"strings"
)

const (
	Initializing = "initializing"
	Running      = "running"
	Pulling      = "pulling"
	Pulled       = "pulled"
	Migrating    = "migrating"
	Cleaning     = "cleaning"
	Borked       = "borked"
	OpFailed     = "opfailed"
	Launching    = "launching"
)

type Version struct {
	Container *docker.Container
	Image     *docker.Image
}
type Service struct {
	Name         string
	State        string
	Current      *Version
	Previous     *Version
	DeployErrror error
}

var Client *docker.Client

func init() {
	var err error
	Client, err = docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
}
func (t *Service) ensureState(state string) {
	if t.State != state {
		panic(fmt.Sprintf("State `%s` has to be `%s`", t.State, state))
	}
}

// the FSM transition map
func (t *Service) transition(state string) (string, error) {
	switch fmt.Sprintf("%s to %s", t.State, state) {
	case fmt.Sprintf("%s to %s", Initializing, Running):
		break
	case fmt.Sprintf("%s to %s", Initializing, Borked):
		break
	case fmt.Sprintf("%s to %s", Running, Pulling):
		break
	case fmt.Sprintf("%s to %s", Running, Launching):
		break
	case fmt.Sprintf("%s to %s", Pulling, OpFailed):
		break
	case fmt.Sprintf("%s to %s", Pulling, Pulled):
		break
	case fmt.Sprintf("%s to %s", Pulled, Launching):
		break
	case fmt.Sprintf("%s to %s", Launching, OpFailed):
		break
	case fmt.Sprintf("%s to %s", Pulled, Migrating):
		break
	case fmt.Sprintf("%s to %s", Migrating, Cleaning):
		break
	case fmt.Sprintf("%s to %s", Cleaning, Running):
		break
	default:
		panic(fmt.Sprintf("Transition from '%s' to '%s' is not allowed", t.State, state))
	}
	t.State = state
	return t.State, nil
}

func CreateService(name string) *Service {
	service := Service{State: Initializing, Name: name}
	return &service
}

/**
* This method is mainly used by testsuite, won't have
* real world uses.
 */
func CreateServiceFromImage(image string) (*Service, error) {
	service := CreateService("Service: " + image)

	containers, err := Client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return nil, err
	}
	var cid string = ""
	var brief *docker.APIContainers
	for _, container := range containers {
		if container.Image == image {
			cid = container.ID
			brief = &container
			break
		}
	}
	if cid == "" {
		service.transition(Borked)
		return service, errors.New(fmt.Sprintf("Container `%s` is not running", image))
	}
	err = service.Latch(cid)
	if err != nil {
		service.transition(Borked)
		return service, err
	}
	if false {
		Println("---------- Brief ---------")
		Println(brief)
		Println("---------- Container ---------")
		Println(service.Current.Container)
		Println("---------- Image ---------")
		Println(service.Current.Image)
	}
	return service, nil
}

/* Latches onto a runing container matched by container.Image
* and populates the container.Current with gathered data
 */
func (t *Service) Latch(cid string) error {
	t.ensureState(Initializing)
	var err error
	ver := Version{}
	ver.Container, err = Client.InspectContainer(cid)
	if err != nil {
		return err
	}

	ver.Image, err = Client.InspectImage(ver.Container.Image)
	if err != nil {
		return err
	}

	if ver.Image == nil {
		return errors.New(fmt.Sprintf("Image `%s` could not be resolved?!?", ver.Container.Image))
	}
	t.Current = &ver

	// Alright we have all the info we need for now, transition to running
	t.transition(Running)
	return nil
}

type imageURL struct {
	Path  string
	Image string
	Tag   string
}

func unpackURL(url string) *imageURL {
	r, _ := regexp.Compile(":([^:]+)$")
	i := imageURL{}
	if m := r.FindStringSubmatch(url); len(m) > 1 {
		i.Tag = m[1]
		url = r.ReplaceAllString(url, "")
	}
	s := strings.Split(url, "/")
	i.Image = s[len(s)-1]
	i.Path = strings.Join(s[0:len(s)-1], "/")
	return &i
}
func (t *imageURL) packURL() string {
	return fmt.Sprintf("%s/%s:%s", t.Path, t.Image, t.Tag)
}
func (t *Version) imageUrl() *imageURL {
	return unpackURL(t.Container.Config.Image)
}

func (t *Version) MigrateConfigTo(tag string) *docker.CreateContainerOptions {
	var hostconfig docker.HostConfig
	var config docker.Config
	config = *t.Container.Config

	// Set new tag to be used
	r, _ := regexp.Compile(":[^:]+$")
	config.Image = r.ReplaceAllString(config.Image, ":"+tag)

	hostconfig = *t.Container.HostConfig

	c := docker.CreateContainerOptions{
		Name:       t.Container.Name,
		Config:     &config,
		HostConfig: &hostconfig,
	}
	return &c
}
func (t *Service) Redeploy(tag string) error {
	t.transition(Pulling)
	url := t.Current.imageUrl()
	url.Tag = tag
	err := Client.PullImage(docker.PullImageOptions{
		Repository: Sprintf("%s/%s", url.Path, url.Image),
		//		Registry:     url.Path,
		Tag:          url.Tag,
		OutputStream: os.Stdout,
	}, docker.AuthConfiguration{})
	if err != nil {
		t.DeployErrror = err
		t.transition(OpFailed)
	} else {
		t.transition(Pulling)
	}

	return err
}
func (t *Service) Healthcheck(version *Version) error {
	// TODO: Implement customizable healthchecks.
	return nil
}
func (t *Service) swapVers() {
	c := t.Current
	t.Current = t.Previous
	t.Previous = c
}
func (t *Service) Migrate(tag string) error {
	t.transition(Launching)
	url := t.Current.imageUrl()
	url.Tag = tag
	creationConf := t.Current.MigrateConfigTo(tag)
	nextImg, err := Client.InspectImage(url.packURL())
	if err != nil {
		t.DeployErrror = err
		t.transition(OpFailed)
		return err
	}
	err = Client.RenameContainer(docker.RenameContainerOptions{
		ID:   t.Current.Container.ID,
		Name: Sprintf("%s_prev", t.Current.Container.Name),
	})

	if err != nil {
		t.DeployErrror = err
		t.transition(OpFailed)
		return err
	}
	nextContainer, err := Client.CreateContainer(*creationConf)
	if err != nil {
		t.DeployErrror = err
		t.transition(OpFailed)
		return err
	}
	t.swapVers()
	t.Current = &Version{Container: nextContainer, Image: nextImg}
	err = Client.StartContainer(nextContainer.ID, docker.HostConfig{})

	if err != nil {
		t.DeployErrror = err
		t.transition(OpFailed)
		return err
	}

	err = t.Healthcheck(t.Current)
	if err != nil {
		t.DeployErrror = err
		t.transition(OpFailed)
		return err
	}

	t.transition(Migrating)

	err = t.updateAndWait()
	if err != nil {
		t.DeployErrror = err
		t.transition(OpFailed)
		return err
	}
	return nil
}

func (t *Service) updateAndWait() error {
	t.ensureState(Migrating)
	// TODO: Put some flag on the t.Previous container
	// so that it gets ignored by the dockergen regenerator.
	changed, err := util.RegenerateConf()
	if !changed || err != nil {
		t.DeployErrror = err
		t.transition(OpFailed)
		return t.DeployErrror
	}
	t.DeployErrror = util.ReloadNginx()
	if t.DeployErrror != nil {
		t.transition(OpFailed)
	} else {
		// TODO: why did i break my previous states?
		// Did i really need an finish-statement?
		// Anyways, this is logically were the cleaning state
		// Should be entererd.
		t.transition(Cleaning)
	}
	return t.DeployErrror
}
