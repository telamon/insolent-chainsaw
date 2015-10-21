package models

import (
	"errors"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	. "github.com/kr/pretty"
	"regexp"
)

const (
	Initializing = "initializing"
	Running      = "running"
	Pulling      = "pulling"
	Migrating    = "migrating"
	Cleaning     = "cleaning"
	Borked       = "borked"
)

type Version struct {
	Container *docker.Container
	Image     *docker.Image
}
type Service struct {
	Name    string
	State   string
	Current *Version
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
	case fmt.Sprintf("%s to %s", Pulling, Migrating):
		break
	case fmt.Sprintf("%s to %s", Migrating, Cleaning):
		break
	case fmt.Sprintf("%s to %s", Cleaning, Running):
		break
	default:
		return t.State, errors.New(fmt.Sprintf("Transition from '%s' to '%s' is not allowed", t.State, state))
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

func (t *Service) Redeploy(tag string) {}
