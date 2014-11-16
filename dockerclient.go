package artifactory

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"regexp"
	"sort"

	"github.com/fsouza/go-dockerclient"
	"github.com/rafecolton/go-dockerclient-sort"
)

var client *docker.Client

// Dockerclient returns the dockerclient used by the artifactory package
func Dockerclient() (*docker.Client, error) {
	if client != nil {
		return client, nil
	}

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

	return client, nil
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

// LatestImageByName uses the provided docker client to get the id
// most-recently-created image with a name matching `name`
func LatestImageByName(client *docker.Client, name string) (string, error) {
	images, err := client.ListImages(false)
	if err != nil {
		return "", err
	}
	sort.Sort(dockersort.ByCreatedDescending(images))
	for _, image := range images {
		for _, tag := range image.RepoTags {
			matched, err := regexp.MatchString("^"+name+"$", tag)
			if err != nil {
				return "", nil
			}
			if matched {
				return image.ID, nil
			}
		}
	}

	return "", fmt.Errorf("unable to find image named %s", name)
}
