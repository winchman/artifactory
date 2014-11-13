package artifactory

import (
	"sync"
)

/*
ResourcePath is a string referring to the path inside the container
to the file or directory.  The type makes it clear what the key
should be when adding it to map
*/
type ResourcePath string

/*
Resource is a type that represents a filepath inside the container
that corresponds to a real file on disk on the host machine.  It has
a nested read-write lock such that it may be locked when being
concurrently read from / written to.
*/
type Resource struct {
	Path    string       // redundant with ResourcePath as the key into the ResourceSet map - only one is needed
	lock    sync.RWMutex // used for reading/writing the state and the actual file
	Present bool
	Error   error
}
