package artifactory

import (
	"io"
	"sync"
)

/*
RWArtifactory is an implementation of the Artifactory interface
*/
type RWArtifactory struct {
	resourceMap  map[Handle]ResourceMap
	sync.RWMutex // for calling Reset() and safety when adding a new ResourceSet
}

/*
Artifactory is a type that can be used to handle all of the artifact-related
interactions for a given build.  It is the responsibility of the caller to
write the resulting artifacts to the correct place on disk once they are
returned
*/
type Artifactory interface {
	// Reset zeros out any data structures that store resource information
	// in memory.  It also deletes the corresponding files from the host
	// filesystem
	Reset() error

	// ResetHandle zeros out the files and data from one given handle
	ResetHandle(Handle) error

	// AddResource gives an Artifactory a list of resource paths, for a
	// given handle, that may be requested by the user.  Nominally, this
	// allows the artifactory to populate the data structure without
	// actually retrieving (and returning) the files from a container.
	//
	// I'm not 100% this function will be necessary.
	AddResource(Handle, ...ResourcePath) error

	// RetrieveResources will return an io.ReadCloser from which the
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
	RetrieveResources(Handle, retrieveAll bool, requestedPaths ...ResourcePath) (map[ResourcePath]io.ReadCloser, error)
}
