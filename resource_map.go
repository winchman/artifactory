package artifactory

import (
	"sync"
)

/*
ResourceMap is a thread-safe hash map of resource values.  Resources may be added
and their artifacts will be available upon request
*/
type ResourceMap struct {
	resources map[ResourcePath]Resource
	lock      sync.RWMutex // for adding new resource objects
}
