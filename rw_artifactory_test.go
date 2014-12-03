package artifactory

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/docker/pkg/archive"
	"github.com/fsouza/go-dockerclient"
	"github.com/rafecolton/go-dockerclient-quick"
)

const (
	imageName    = "quay.io/rafecolton/docker-builder:latest"
	resourcePath = "/app/bin/docker-builder"
)

func TestArtifactory(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration test")
	}

	tempDir, err := ioutil.TempDir("", "artifactory-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	var art = NewArtifactory(tempDir)
	var image *docker.APIImages
	client, err := dockerclient.NewDockerClient()
	if err != nil {
		t.Fatal(err)
	}
	image, err = client.LatestImageByRegex("^" + imageName + "$")
	if err != nil {
		t.Fatal(err)
	}

	if err := art.AddResource(image.ID, resourcePath); err != nil {
		t.Fatal(err)
	}
	var resourceFunc = func(r *Resource, err error) error {
		if err != nil {
			return err
		}
		resourceBytes, err := r.ArtifactBytes()
		if err != nil {
			t.Fatal(err)
		}
		byteReader := bytes.NewReader(resourceBytes)
		archive.Untar(byteReader, os.Getenv("PWD"), &archive.TarOptions{})
		return nil
	}

	if err := art.EachResource(image.ID, resourceFunc); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(os.Getenv("PWD") + "/docker-builder"); os.IsNotExist(err) {
		t.Fatal(err)
	}

	if err = os.RemoveAll(os.Getenv("PWD") + "/docker-builder"); err != nil {
		t.Fatal(err)
	}
}
