package artifactory

import (
	"sync"
)

// NewArtifactory produces an initialized instance of a struct that implements
// the Artifactory interface
func NewArtifactory(storageDir string) Artifactory {
	return &RWArtifactory{
		resourceMap: map[string]*ResourceSet{},
		storageDir:  storageDir,
	}
}

/*
RWArtifactory is an implementation of the Artifactory interface
*/
type RWArtifactory struct {
	resourceMap map[string]*ResourceSet
	storageDir  string
	lock        sync.RWMutex // for calling Reset() and safety when adding a new ResourceSet
}

// Reset zeros out any data structures that store resource information
// in memory.  It also deletes the corresponding files from the host
// filesystem
func (art *RWArtifactory) Reset() error {
	for handle := range art.resourceMap {
		if err := art.ResetHandle(handle); err != nil {
			return err
		}
	}
	return nil
}

// ResetHandle zeros out the files and data from one given handle
func (art *RWArtifactory) ResetHandle(h string) error {
	art.lock.Lock()
	defer art.lock.Unlock()
	set := art.resourceMap[h]
	if set == nil {
		return nil
	}
	set.Each(func(r *Resource, err error) error { return r.Reset() })
	art.resourceMap[h] = nil

	return nil
}

// AddResource gives an Artifactory a list of resource paths, for a
// given handle, that may be requested by the user.  Nominally, this
// allows the artifactory to populate the data structure without
// actually retrieving (and returning) the files from a container.
func (art *RWArtifactory) AddResource(h string, resourcePaths ...string) error {
	art.lock.Lock()
	defer art.lock.Unlock()
	if art.resourceMap[h] == nil {
		art.resourceMap[h] = NewResourceSet()
	}

	for _, resourcePath := range resourcePaths {
		err := art.resourceMap[h].Add(NewResource(NewResourceOptions{
			Path:       resourcePath,
			StorageDir: art.storageDir + "/" + h,
			Handle:     h,
		}))
		if err != nil {
			return err
		}
	}

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
func (art *RWArtifactory) EachResource(h string, resourceFunc func(*Resource, error) error) error {
	art.lock.RLock()
	defer art.lock.RUnlock()
	set := art.resourceMap[h]
	return set.Each(resourceFunc)
}
