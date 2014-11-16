package artifactory

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestArtifactBytes(t *testing.T) {
	testingDir := os.Getenv("PWD") + "/_testing"
	opts := NewResourceOptions{StorageDir: testingDir, Path: "/app/bin", test: true}
	var resource = NewResource(opts)

	expectedBytes, err := ioutil.ReadFile(testingDir + "/4fc03566a597e6e389ef3c567447520705d78efb.tar")
	if err != nil {
		t.Error(err)
	}
	bytes, err := resource.ArtifactBytes()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(bytes, expectedBytes) {
		t.Errorf("got wrong artifact bytes")
	}
}
