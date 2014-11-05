package artifactory

import (
	"testing"
)

func Test_HandleCreation(t *testing.T) {
	var containerID = "abc123def456"
	var expected = "90bd1b48e958257948487b90bee080ba5ed00caa"
	handle := NewHandle(containerID)
	if handle.String() != expected {
		t.Errorf("expected %q, got %q", expected, handle)
	}
}
