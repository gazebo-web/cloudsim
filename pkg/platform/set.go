package platform

import "errors"

var (
	// ErrPlatformNotFound is returned when a platform for a provided selector cannot be found.
	ErrPlatformNotFound = errors.New("platform not found")
)

// Selector is used to uniquely identify a Platform.
type Selector string

// Manager manages a set of platforms that can be used by Cloudsim.
// Implementations must be simple platform containers and must not be aware of the feature sets, differences (if any) or
// implementation details of the available platforms they contain. This is by design, in order to give applications the
// flexibility to manage sets of platforms as they see fit.
// All contained platforms are completely independent from each other. As an example, you can have a set of platforms
// that use the same internal components (e.g. AWS+Kubernetes), but that are configured to point to different
// service regions. At the same time, it is possible to have additional platforms making use of entirely different
// components (e.g. GCP+Mesos, Azure+Swarm) co-exist with the set of AWS+Kubernetes platforms.
// Every platform is uniquely identified by a Selector, a user-defined identifier. To make use of a specific platform,
// the target platform's Selector is passed to the Manager implementation through the Platform method.
// An example of a selector set:
//   * aws_k8s_us_east_1 - Platform containing AWS and Kubernetes components pointed at us-east-1.
//   * aws_k8s_us_east_2 - Platform containing AWS and Kubernetes components pointed at us-east-2.
//   * gcp_mesos - Platform containing GCE and Apache Mesos components.
//   * azure_swarm - Platform containing Azure and Docker Swarm components.
type Manager interface {
	// Selectors returns a slice with all the available platform selectors.
	Selectors() []Selector
	// Platforms returns a slice with all the available platforms.
	Platforms() []Platform
	// Platform returns the platform that matches a specific selector, or an error if it is not found.
	Platform(selector Selector) (Platform, error)
}

// Map is the default Manager implementation.
type Map map[Selector]Platform

// Selectors returns a slice with all the available platform selectors.
func (p Map) Selectors() []Selector {
	selectors := make([]Selector, len(p))
	i := 0
	for selector := range p {
		selectors[i] = selector
		i++
	}

	return selectors
}

// Platforms returns a slice with all the available platforms.
func (p Map) Platforms() []Platform {
	platforms := make([]Platform, len(p))
	i := 0
	for _, platform := range p {
		platforms[i] = platform
		i++
	}

	return platforms
}

// Platform returns receives a selector and returns its matching platform or an error if it is not found.
func (p Map) Platform(selector Selector) (Platform, error) {
	platform, ok := p[selector]
	if !ok {
		return nil, ErrPlatformNotFound
	}

	return platform, nil
}
