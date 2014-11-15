package artifactory

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path"

	"github.com/fsouza/go-dockerclient"
)

func (r *RWResource) checkAndPopulate() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if !r.present {
		client, err := dockerclient()
		if err != nil {
			return err
		}

		containerID, err := createAndStartContainer(client, r.handle.String())
		if err != nil {
			return err
		}

		defer func() {
			go killContainer(client, containerID)
		}()

		var buf bytes.Buffer

		opts := docker.CopyFromContainerOptions{
			Container:    containerID,
			Resource:     string(r.path),
			OutputStream: &buf,
		}
		if err := client.CopyFromContainer(opts); err != nil {
			return err
		}
		if err := ioutil.WriteFile(r.storageDir+"/"+r.artifactFileName(), buf.Bytes(), 0644); err != nil {
			return err
		}

		r.present = true
	}
	return nil
}

func dockerclient() (client *docker.Client, err error) {
	endpoint, err := getEndpoint()
	if err != nil {
		return nil, err
	}
	certPath := os.Getenv("DOCKER_CERT_PATH")
	tlsVerify := os.Getenv("DOCKER_TLS_VERIFY") != ""

	if endpoint.Scheme == "https" {
		cert := path.Join(certPath, "cert.pem")
		key := path.Join(certPath, "key.pem")
		ca := ""
		if tlsVerify {
			ca = path.Join(certPath, "ca.pem")
		}

		client, err = docker.NewTLSClient(endpoint.String(), cert, key, ca)
		if err != nil {
			return nil, err
		}
	} else {
		client, err = docker.NewClient(endpoint.String())
		if err != nil {
			return nil, err
		}
	}

	return
}

func getEndpoint() (*url.URL, error) {
	endpoint := os.Getenv("DOCKER_HOST")
	if endpoint == "" {
		endpoint = "unix:///var/run/docker.sock"
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse endpoint %s as URL", endpoint)
	}
	if u.Scheme == "tcp" {
		_, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s for port", u.Host)
		}

		// Only reliable way to determine if we should be using HTTPS appears to be via port
		if os.Getenv("DOCKER_HOST_SCHEME") != "" {
			u.Scheme = os.Getenv("DOCKER_HOST_SCHEME")
		} else if port == "2376" {
			u.Scheme = "https"
		} else {
			u.Scheme = "http"
		}
	}
	return u, nil
}

// creates the container with image id and starts it, returns the container id
func createAndStartContainer(client *docker.Client, id string) (string, error) {
	createOpts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:      id,
			Cmd:        []string{"60"},
			Entrypoint: []string{"sleep"},
		},
	}

	container, err := client.CreateContainer(createOpts)
	if err != nil {
		fmt.Println("error creating container: " + err.Error())
		return "", err
	}

	fmt.Println("starting container for artifact extraction...")
	if err := client.StartContainer(container.ID, &docker.HostConfig{}); err != nil {
		fmt.Println("error starting container: " + err.Error())
		return "", err
	}
	return container.ID, nil
}

// kills the container
func killContainer(client *docker.Client, containerID string) {
	fmt.Println("artifact extraction complete, killing container")
	opts := docker.KillContainerOptions{
		ID:     containerID,
		Signal: docker.SIGKILL,
	}
	if err := client.KillContainer(opts); err != nil {
		fmt.Println("error killing container: " + err.Error())
	}
}
