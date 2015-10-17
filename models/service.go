package models

import (
	"errors"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	Initializing = "initializing"
	Running      = "running"
	Pulling      = "pulling"
	Migrating    = "migrating"
	Cleaning     = "cleaning"
)

type Version struct {
	CBrief    *docker.APIContainers
	Container *docker.Container
	Image     *docker.APIImages
}
type Service struct {
	Image   string
	State   string
	Current Version
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
		panic(fmt.Sprintf("Expected state `%s` to be `%s`", t.State, state))
	}
}

// the FSM transition map
func (t *Service) transition(state string) (string, error) {
	switch fmt.Sprintf("%s to %s", t.State, state) {
	case fmt.Sprintf("%s to %s", Initializing, Running):
		break
	case fmt.Sprintf("%s to %s", Initializing, Pulling):
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

func CreateService(image string) (*Service, error) {
	service := Service{State: Initializing, Image: image}
	return &service, nil
}

func (t *Service) Latch() error {
	t.ensureState(Initializing)
	containers, err := Client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return err
	}
	t.Current = Version{}
	for _, container := range containers {
		if container.Image == t.Image {
			t.Current.CBrief = &container
			break
		}
	}
	if t.Current.CBrief == nil {
		return errors.New(fmt.Sprintf("Container `%s` is not running", t.Image))
	}
	t.Current.Container, err = Client.InspectContainer(t.Current.CBrief.ID)
	if err != nil {
		return err
	}

	images, err := Client.ListImages(docker.ListImagesOptions{})
	if err != nil {
		return err
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == t.Current.CBrief.Image || tag == t.Current.CBrief.Image+":latest" {
				t.Current.Image = &image
				break
			}
		}
		if t.Current.Image != nil {
			break
		}
	}
	if t.Current.Image == nil {
		return errors.New(fmt.Sprintf("Image `%s` could not be resolved", t.Image))
	}

	// Alright we have all the info we need for now, transition to running
	t.transition(Running)
	return nil
}
