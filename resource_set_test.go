package artifactory

import (
	"testing"
)

func TestResourceSetAddAndGet(t *testing.T) {
	var r = NewResourceSet()
	var resource = NewResource(NewResourceOptions{Path: "foo", test: true})

	err := r.Add(resource)
	if err != nil {
		t.Errorf("adding to resource set failed")
	}

	ret := r.Get("foo")
	if ret == nil || ret.Path() != "foo" {
		t.Errorf("resource should be retrievable")
	}
}

func TestResourceSetAddTwice(t *testing.T) {
	var r = NewResourceSet()
	var resource = NewResource(NewResourceOptions{Path: "foo", test: true})

	_ = r.Add(resource)
	err := r.Add(resource)
	if err == nil || !IsAlreadyPresentInSetError(err) {
		t.Errorf("adding the same resource twice should error")
	}
}
