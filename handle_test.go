package artifactory

import (
	"testing"
)

func Test_HandleCreation(t *testing.T) {
	var containerID = "abc123def456"
	var expected = "abc123def456"
	handle := NewHandle(containerID)
	if handle.String() != expected {
		t.Errorf("expected %q, got %q", expected, handle)
	}
}
