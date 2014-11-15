package artifactory

import (
	"sync"
)

// NewArtifactory produces an initialized instance of a struct that implements
// the Artifactory interface
func NewArtifactory(storageDir string) Artifactory {
	return &RWArtifactory{
		resourceMap: map[Handle]*ResourceSet{},
		storageDir:  storageDir,
	}
}

/*
RWArtifactory is an implementation of the Artifactory interface
*/
type RWArtifactory struct {
	resourceMap  map[Handle]*ResourceSet
	storageDir   string
	sync.RWMutex // for calling Reset() and safety when adding a new ResourceSet
}

// Reset zeros out any data structures that store resource information
// in memory.  It also deletes the corresponding files from the host
// filesystem
func (art *RWArtifactory) Reset() error {
	return nil
}

// ResetHandle zeros out the files and data from one given handle
func (art *RWArtifactory) ResetHandle(Handle) error {
	return nil
}

// AddResource gives an Artifactory a list of resource paths, for a
// given handle, that may be requested by the user.  Nominally, this
// allows the artifactory to populate the data structure without
// actually retrieving (and returning) the files from a container.
//
// I'm not 100% this function will be necessary.
func (art *RWArtifactory) AddResource(Handle, ...ResourcePath) error {
	return nil
}

// EachResource will return an io.ReadCloser from which the
// file contents can be read for each resource.  The file contents
// for each will be a tarball (compressed?) such that it can be
// passed directly into the docker `archive` package's
// DecompressSteam or Untar function. The intent is that the
// resource be untarred / decompressed into
// `$CONTEXT_DIR/$PREFIX/$RESOURCE_PATH` where $CONTEXT_DIR is the
// directory from which the dependent image will be built, $PREFIX
// is an arbitrary prefix (e.g. "inbox"), and $RESOURCE_PATH is the
// full path at which the resource can be found *inside* the
// container
func (art *RWArtifactory) EachResource(Handle, func(*Resource, error)) error {
	return nil
}
