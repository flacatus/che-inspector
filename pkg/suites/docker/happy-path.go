package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/flacatus/che-inspector/pkg/common/instance"
	"io"
	"os"
)

var amd []string

func RunHappyPathDocker(dockerClient *client.Client, testSpec *instance.CheTestsSpec) (err error) {
	if err := PullHappyPathImage(dockerClient, testSpec); err != nil {
		return err
	}

	if err := CreateAndStartContainer(dockerClient, testSpec); err != nil {
		return err
	}

	return nil
}

func PullHappyPathImage(dockerClient *client.Client, testSpec *instance.CheTestsSpec) (err error) {
	reader, err := dockerClient.ImagePull(context.Background(), testSpec.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	_, _ = io.Copy(os.Stdout, reader)
	return nil
}

func CreateAndStartContainer(dockerClient *client.Client, testSpec *instance.CheTestsSpec) (err error) {
	for _, e := range testSpec.Env {
		amd = append(amd, e.Name + "=" + e.Value)
	}
	resp, err := dockerClient.ContainerCreate(context.Background(), &container.Config{
		Image: testSpec.Image,
		Env: amd,
		Tty:   false,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type: mount.TypeBind,
				Source: testSpec.Artifacts.To,
				Target: testSpec.Artifacts.FromContainerPath,
			},
		},
	}, nil, testSpec.Name)
	if err != nil {
		return err
	}
	if err := dockerClient.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	out, err := dockerClient.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{ShowStdout: true, Follow: true})
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, out)

	return err
}
