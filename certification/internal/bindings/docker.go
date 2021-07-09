package bindings

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	specsv1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func (d *DockerClient) PullImage(nameOrID string) error {
	reader, err := d.Client.ImagePull(d.Context, nameOrID, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	return nil
}

func (d *DockerClient) SaveImage(nameOrID string) error {
	outfile, err := os.Create("image.tar.gz")
	if err != nil {
		return err
	}
	defer outfile.Close()

	reader, err := d.Client.ImageSave(d.Context, []string{nameOrID})
	if err != nil {
		return err
	}
	defer reader.Close()

	outfile.ReadFrom(reader)
	return nil
}

func (d *DockerClient) ListImages() ([]Image, error) {
	list, err := d.Client.ImageList(d.Context, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}

	var images []Image
	for _, image := range list {
		images = append(images, Image{
			ID:          image.ID,
			Labels:      image.Labels,
			RepoTags:    image.RepoTags,
			RepoDigests: image.RepoDigests,
		})
	}

	return images, nil
}

func (d *DockerClient) RemoveImage(nameOrID string) ([]string, error) {
	responses, err := d.Client.ImageRemove(d.Context, nameOrID, types.ImageRemoveOptions{})
	if err != nil {
		return nil, err
	}

	var deleted []string
	for _, response := range responses {
		deleted = append(deleted, response.Deleted)
	}

	return deleted, nil
}

func (d *DockerClient) RunContainer(nameOrID string, containerName string) (string, error) {
	if err := d.PullImage(nameOrID); err != nil {
		return "", err
	}

	options := &container.Config{
		Image: nameOrID,
		Cmd:   []string{},
		Tty:   false,
	}
	hostConfig := &container.HostConfig{}
	netConfig := &network.NetworkingConfig{}
	platform := &specsv1.Platform{}
	container, err := d.Client.ContainerCreate(d.Context, options, hostConfig, netConfig, platform, containerName)
	if err != nil {
		return "", err
	}

	return container.ID, d.Client.ContainerStart(d.Context, container.ID, types.ContainerStartOptions{})
}

func (d *DockerClient) ListContainers() ([]Container, error) {
	list, err := d.Client.ContainerList(d.Context, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	var containers []Container
	for _, container := range list {
		containers = append(containers, Container{
			ID:     container.ID,
			Names:  container.Names,
			Labels: container.Labels,
		})
	}

	return containers, nil
}

func (d *DockerClient) RemoveContainer(nameOrID string) error {
	return d.Client.ContainerRemove(d.Context, nameOrID, types.ContainerRemoveOptions{})
}
