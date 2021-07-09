package bindings

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/docker/docker/client"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/certification/errors"
)

func podmanSocket() string {
	socketPath := []string{
		os.Getenv("XDG_RUNTIME_DIR"),
		"podman",
		"podman.sock",
	}
	if _, err := os.Stat(filepath.Join(socketPath...)); err != nil {
		socketPath = []string{
			"/run",
			"podman",
			"podman.sock",
		}
	}

	return fmt.Sprintf("unix://%s", filepath.Join(socketPath...))
}

func Discovery() (*DiscoveryClient, error) {
	discovery := &DiscoveryClient{}
	ctx, _ := bindings.NewConnection(context.Background(), podmanSocket())
	if ctx != nil {
		discovery.podman = &PodmanClient{
			Context: ctx,
		}
	}

	// check if docker socket exists
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if cli != nil {
		discovery.docker = &DockerClient{
			Client:  cli,
			Context: context.Background(),
		}
	}

	if !discovery.dockerExists() && !discovery.podmanExists() {
		return nil, errors.ErrNoSocketFound
	}

	return discovery, nil
}

func (d *DiscoveryClient) Client() ContainerTool {
	if d.podmanExists() {
		return d.podman
	}

	return d.docker
}

func (d *DiscoveryClient) podmanExists() bool {
	return d.podman != nil
}

func (d *DiscoveryClient) dockerExists() bool {
	return d.docker != nil
}
