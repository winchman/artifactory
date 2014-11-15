package artifactory

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
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
type Resource interface {
	ArtifactBytes() ([]byte, error)
	Path() ResourcePath
}

type RWResource struct {
	Error error

	handle     Handle
	lock       sync.RWMutex // used for reading/writing the state and the actual file
	path       ResourcePath
	present    bool
	storageDir string
}

/*
NewResourceOptions is a struct to disambiguate the options passed to
NewResource.  Handle should correspond to a valid containerID as created
by NewHandle to ensure that artifact extraction is possible.  If being used
for testing, Handle may be nil, and resource.present should be set to true.
*/
type NewResourceOptions struct {
	StorageDir string
	Handle     Handle
	Path       string
	test       bool // private, can only be set for tests in same package
}

// Generally, this would be the storage dir of the artifactory plus the
// uniqueID of the handle
func NewResource(opts NewResourceOptions) Resource {
	return &RWResource{
		storageDir: opts.StorageDir,
		path:       ResourcePath(opts.Path),
		present:    opts.test, // if testing mode, mark as already present
		handle:     opts.Handle,
	}
}

func (r *RWResource) ArtifactBytes() ([]byte, error) {
	if err := r.checkAndPopulate(); err != nil {
		return nil, err
	}
	return r.artifactBytes()
}

func (r *RWResource) artifactBytes() ([]byte, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return ioutil.ReadFile(r.storageDir + "/" + r.artifactFileName())
}

func (r *RWResource) Path() ResourcePath {
	return r.path
}

func (r *RWResource) artifactFileName() string {
	var hasher = sha1.New()
	io.WriteString(hasher, string(r.path))
	return fmt.Sprintf("%x", hasher.Sum(nil)) + ".tar"
}
