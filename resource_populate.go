package artifactory

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fsouza/go-dockerclient"
)

func (r *Resource) checkAndPopulate() error {
	if !r.present {
		r.lock.Lock()
		defer r.lock.Unlock()
		client, err := NewDockerClient()
		if err != nil {
			return err
		}

		containerID, err := client.createAndStartContainer(r.handle.String())
		if err != nil {
			return err
		}

		defer func() {
			go client.killContainer(containerID)
		}()

		var buf bytes.Buffer

		opts := docker.CopyFromContainerOptions{
			Container:    containerID,
			Resource:     string(r.path),
			OutputStream: &buf,
		}
		if err := (*docker.Client)(client).CopyFromContainer(opts); err != nil {
			return err
		}
		if err := os.MkdirAll(r.storageDir, 0777); err != nil {
			return err
		}
		if err := ioutil.WriteFile(r.storageDir+"/"+r.artifactFileName(), buf.Bytes(), 0644); err != nil {
			return err
		}

		r.present = true
	}
	return nil
}

// creates the container with image id and starts it, returns the container id
func (client *DockerClient) createAndStartContainer(id string) (string, error) {
	createOpts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:      id,
			Cmd:        []string{"60"},
			Entrypoint: []string{"sleep"},
		},
	}

	container, err := (*docker.Client)(client).CreateContainer(createOpts)
	if err != nil {
		//fmt.Println("error creating container: " + err.Error())
		return "", err
	}

	//fmt.Println("starting container for artifact extraction...")
	if err := (*docker.Client)(client).StartContainer(container.ID, &docker.HostConfig{}); err != nil {
		//fmt.Println("error starting container: " + err.Error())
		return "", err
	}
	return container.ID, nil
}

// kills the container
func (client *DockerClient) killContainer(containerID string) {
	//fmt.Println("artifact extraction complete, killing container")
	opts := docker.KillContainerOptions{
		ID:     containerID,
		Signal: docker.SIGKILL,
	}
	if err := (*docker.Client)(client).KillContainer(opts); err != nil {
		fmt.Println("error killing container: " + err.Error())
	}
}
