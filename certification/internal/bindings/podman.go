package bindings

import (
	"os"

	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/containers/podman/v3/pkg/specgen"
)

func (p *PodmanClient) PullImage(nameOrID string) error {
	_, err := images.Pull(p.Context, nameOrID, &images.PullOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (p *PodmanClient) SaveImage(nameOrID string) error {
	outfile, err := os.Create("image.tar.gz")
	if err != nil {
		return err
	}
	defer outfile.Close()

	var compress bool = true
	err = images.Export(p.Context, []string{nameOrID}, outfile, &images.ExportOptions{
		Compress: &compress,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *PodmanClient) ListImages() ([]Image, error) {
	list, err := images.List(p.Context, &images.ListOptions{})
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

func (p *PodmanClient) RemoveImage(nameOrID string) ([]string, error) {
	responses, err := images.Remove(p.Context, []string{nameOrID}, &images.RemoveOptions{})
	if err != nil {
		return nil, err[0]
	}

	return responses.Deleted, nil
}

func (p *PodmanClient) RunContainer(nameOrID string, containerName string) (string, error) {
	if err := p.PullImage(nameOrID); err != nil {
		return "", err
	}

	containerSpec := &specgen.SpecGenerator{
		ContainerBasicConfig: specgen.ContainerBasicConfig{
			Name:       containerName,
			Entrypoint: []string{},
			Command:    []string{},
		},
	}
	container, err := containers.CreateWithSpec(p.Context, containerSpec, &containers.CreateOptions{})
	if err != nil {
		return "", err
	}

	return container.ID, containers.Start(p.Context, container.ID, &containers.StartOptions{})
}

func (d *PodmanClient) ListContainers() ([]Container, error) {
	list, err := containers.List(d.Context, &containers.ListOptions{})
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

func (p *PodmanClient) RemoveContainer(nameOrID string) error {
	var forceDelete bool = true
	options := &containers.RemoveOptions{
		Force: &forceDelete,
	}
	return containers.Remove(p.Context, nameOrID, options)
}
