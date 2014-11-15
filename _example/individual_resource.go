package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/docker/docker/pkg/archive"
	"github.com/fsouza/go-dockerclient"
	"github.com/rafecolton/go-dockerclient-sort"
	"github.com/sylphon/artifactory"
)

const (
	imageName    = "quay.io/rafecolton/docker-builder:latest"
	resourcePath = "/app/bin/docker-builder"
)

// The example details how to extract the docker-builder binary from the lastest
// docker-builder image
func main() {
	client, err := artifactory.Dockerclient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	containerID, err := latestImageByName(client, imageName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	opts := artifactory.NewResourceOptions{
		StorageDir: os.Getenv("PWD"),
		Handle:     artifactory.NewHandle(containerID),
		Path:       resourcePath,
	}

	var resource = artifactory.NewResource(opts)

	artifactBytes, err := resource.ArtifactBytes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	byteReader := bytes.NewReader(artifactBytes)
	archive.Untar(byteReader, os.Getenv("PWD"), &archive.TarOptions{})
	_ = resource.Reset()
}

func latestImageByName(client *docker.Client, name string) (string, error) {
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