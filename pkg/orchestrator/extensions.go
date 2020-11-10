package orchestrator

// Extensions is a generic interface to represent an extension of the cluster behavior.
// We use the Extensions interface to represent very specific use cases of our cluster, any Extensions implementation
// should be created if it doesn't make sense to include that behavior on one of the existent components.
// As soon as it makes sense for the implementation be considered part of a specific component, a new component should be
// created and the extension should be marked as nil, and all calls should be replaced with the specific component.
type Extensions interface{}
