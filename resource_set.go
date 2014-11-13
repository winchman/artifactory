package artifactory

import (
	"sync"
)

/*
ResourceSet is a thread-safe hash map of resource values.  Resources may be added
and their artifacts will be available upon request
*/
type ResourceSet struct {
	resources map[ResourcePath]Resource
	lock      sync.RWMutex // for adding new resource objects
}
