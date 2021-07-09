package bindings

import (
	"context"

	"github.com/docker/docker/client"
)

// PodmanClient is a client used to interact with
// containers and images via the with podman socket
type PodmanClient struct {
	context.Context
}

// DockerSocket is a client used to interact with
// containers and images via the with podman socket
type DockerClient struct {
	*client.Client
	context.Context
}

// DiscoveryClient is a custom client used when discovering
// whether podman or docker sockets are available for use by
// the preflight tool
type DiscoveryClient struct {
	podman *PodmanClient
	docker *DockerClient
}

// Container struct provides a way to standardize the results returned
// by the methods run by the podman client and docker client methods
type Container struct {
	ID     string
	Names  []string
	Labels map[string]string
}

// Image struct provides a way to standardize the results returned
// by the methods run by the podman client and docker client methods
type Image struct {
	ID          string
	Labels      map[string]string
	RepoTags    []string
	RepoDigests []string
}

// The ContainerTool interface provides a unified way to interact with
// both the podman and docker methods below
type ContainerTool interface {
	// PullImage takes an image name and pulls the image to the local cache
	PullImage(nameOrID string) error
	// ListImages returns a list of images available on the local cache
	ListImages() ([]Image, error)
	// RemoveImage deletes an image from the local cache
	RemoveImage(nameOrID string) ([]string, error)
	// RunContainer creates and starts a container
	RunContainer(nameOrID string, containerName string) (string, error)
	// ListContainer lists the containers on the local system
	ListContainers() ([]Container, error)
	// RemoveContainer deletes a container running on the local system
	RemoveContainer(nameOrID string) error
}
