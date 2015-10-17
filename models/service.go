package models

import (
	"errors"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	. "github.com/kr/pretty"
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
	Brief     *docker.APIContainers
	Container *docker.Container
	Image     *docker.Image
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

func CreateService(image string) (*Service, error) {
	service := Service{State: Initializing, Image: image}
	return &service, nil
}

/* Latches onto a runing container matched by container.Image
* and populates the container.Current with gathered data
 */
func (t *Service) Latch() error {
	t.ensureState(Initializing)
	containers, err := Client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return err
	}
	t.Current = Version{}
	for _, container := range containers {
		if container.Image == t.Image {
			t.Current.Brief = &container
			break
		}
	}
	if t.Current.Brief == nil {
		return errors.New(fmt.Sprintf("Container `%s` is not running", t.Image))
	}
	Println("---------- Brief ---------")
	Println(t.Current.Brief)
	t.Current.Container, err = Client.InspectContainer(t.Current.Brief.ID)
	if err != nil {
		return err
	}

	Println("---------- Container ---------")
	Println(t.Current.Container)
	t.Current.Image, err = Client.InspectImage(t.Current.Container.Image)
	if err != nil {
		return err
	}

	if t.Current.Image == nil {
		return errors.New(fmt.Sprintf("Image `%s` could not be resolved?!?", t.Image))
	}
	Println("---------- Image ---------")
	Println(t.Current.Image)
	// Alright we have all the info we need for now, transition to running
	t.transition(Running)
	return nil
}
