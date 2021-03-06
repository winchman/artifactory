package artifactory

import (
	"errors"
	"sync"
)

// AlreadyPresentInSetErrorMessage is the message returned for AlreadyPresentInSetError errors
const AlreadyPresentInSetErrorMessage = "resource with given path already present in set"

/*
ResourceSet is a thread-safe hash map of resource values.  Resources may be added
and their artifacts will be available upon request
*/
type ResourceSet struct {
	resources map[string]*Resource
	lock      sync.RWMutex // for adding new resource objects
}

// NewResourceSet creates a fully initialized ResourceSet
func NewResourceSet() *ResourceSet {
	return &ResourceSet{
		resources: map[string]*Resource{},
	}
}

/*
Add adds resource r to the set. If a resource is already present with the same
path (r.Path), Add will return an IsIsPresentInSetError
*/
func (set *ResourceSet) Add(r *Resource) error {
	set.lock.Lock()
	defer set.lock.Unlock()

	if set.resources[r.Path()] == nil {
		set.resources[r.Path()] = r
	} else {
		return newAlreadyPresentInSetError()
	}

	return nil
}

// Get returns the resource that exists at path (or nil)
func (set *ResourceSet) Get(path string) *Resource {
	set.lock.RLock()
	defer set.lock.RUnlock()
	return set.resources[path]
}

// AlreadyPresentInSetError is returned on an unsuccessful call to Add()
type AlreadyPresentInSetError error

func newAlreadyPresentInSetError() AlreadyPresentInSetError {
	return errors.New(AlreadyPresentInSetErrorMessage)
}

// IsAlreadyPresentInSetError checks if error e is of the type AlreadyPresentInSetError
func IsAlreadyPresentInSetError(e error) bool {
	if e == nil {
		return false
	}
	switch e.(type) {
	case AlreadyPresentInSetError:
		return true
	default:
		return e.Error() == AlreadyPresentInSetErrorMessage
	}

}

// Each iterates over each resource in the set in a threadsafe manner, yielding
// each one to function resourceFunc.  If resourceFunc returns an error, that
// error is passed to the subseqwuent invocation.
func (set *ResourceSet) Each(resourceFunc func(r *Resource, error error) error) error {
	set.lock.RLock()
	defer set.lock.RUnlock()
	var err error
	for _, resource := range set.resources {
		err = resourceFunc(resource, err)
	}
	return nil
}
