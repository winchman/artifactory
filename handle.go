package artifactory

// Handle is a unique identifier that an Artifactory uses to differentiate between artifacts
type Handle string

/*
NewHandle returns the handle used by the calling artifactory to reference a
given container.  The handle will be a unique identifier, and it will be passed
to the artifactory on requests to operate on artifacts.

Current, containerID is simply returned, as container IDs are unique
*/
func NewHandle(containerID string) Handle {
	return Handle(containerID)
}

// String returns the handle as a string
func (h Handle) String() string {
	return string(h)
}
