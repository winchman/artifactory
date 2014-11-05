package artifactory

//[>
//TODO:
//- implementation
//- tests (including linting and fmtpolice)
//- comments for the linter
//- some way of indicating that the artifactory is preparing a given artifact
//*/

//type Artifactory struct {
//storageDir string
//prepared   map[Handle]bool
//artifacts  map[Handle][]artifact //tentative
//}

////tentative
//type artifact struct {
//resource string
//status   string // i.e. currently preparing so don't try to prepare, just wait
//}

//func NewArtifactory(storageDir string) *Artifactory {
//// validation of storage dir?
//return &Artifactory{storageDir: storageDir}
//}

//type PrepareArtifactOptions struct {
//Handle   Handle
//Resource string // path to file or directory inside the container
//Force    bool   // by default, do not "prepare" artifacts that have already been retrieved
//}

//// FIXME: should this be a private method?  guess it depends on whether we want eager or lazy loading
//// FIXME: do we even want to store to disk at all?  might be more efficient to just stream the tarball right from the container to the requesting function
//func (art *Artifactory) PrepareArtifact(opts PrepareArtifactOptions) error {
//// TODO: normalize the resource name
//// TODO: some way to indicate that the resource is currently being prepared so we don't block on it unless necessary
//return nil
//}

//type PopulateArtifactOptions struct {
//Handle      Handle
//Resource    string // do we want this here?
//PopulateAll bool   // maybe if this is set to true we don't specify a specific resource?
//Destination string
//}

//func (art *Artifactory) PopulateArtifacts(opts PopulateArtifactOptions) error {
//return nil
//}
