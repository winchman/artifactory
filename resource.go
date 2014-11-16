package artifactory

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
	Error error

	handle     Handle
	lock       sync.RWMutex // used for reading/writing the state and the actual file
	path       ResourcePath
	present    bool
	storageDir string
}

// Reset deletes the underlying extracted archive file and resets the state of
// the resource
func (r *Resource) Reset() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.present {
		r.present = false
		if err := os.RemoveAll(r.artifactFullPath()); err != nil {
			return err
		}
	}
	return nil
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

// NewResource returns a properly initialized resource
func NewResource(opts NewResourceOptions) *Resource {
	return &Resource{
		storageDir: opts.StorageDir,
		path:       ResourcePath(opts.Path),
		present:    opts.test, // if testing mode, mark as already present
		handle:     opts.Handle,
	}
}

// ArtifactBytes returns the bytes of the artifact (a `.tar` archive)
func (r *Resource) ArtifactBytes() ([]byte, error) {
	if err := r.checkAndPopulate(); err != nil {
		return nil, err
	}
	return r.artifactBytes()
}

func (r *Resource) artifactBytes() ([]byte, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return ioutil.ReadFile(r.artifactFullPath())
}

// Path returns, for the given resource, the path inside the container at which
// it can be found - used as a unique index for a given handle (container ID)
func (r *Resource) Path() ResourcePath {
	return r.path
}

func (r *Resource) artifactFileName() string {
	var hasher = sha1.New()
	io.WriteString(hasher, string(r.path))
	return fmt.Sprintf("%x", hasher.Sum(nil)) + ".tar"
}

func (r *Resource) artifactFullPath() string {
	return r.storageDir + "/" + r.artifactFileName()
}
