package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/docker/docker/pkg/archive"
	"github.com/rafecolton/go-dockerclient-quick"
	"github.com/sylphon/artifactory"
)

const (
	imageName    = "quay.io/rafecolton/docker-builder:latest"
	resourcePath = "/app/bin/docker-builder"
)

// The example details how to extract the docker-builder binary from the lastest
// docker-builder image
func main() {
	client, err := dockerclient.NewDockerClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	containerID, err := client.LatestImageIDByName(imageName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	opts := artifactory.NewResourceOptions{
		StorageDir: os.Getenv("PWD"),
		Handle:     containerID,
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
