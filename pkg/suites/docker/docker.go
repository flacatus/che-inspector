package docker

import (
	"context"
	"io"
	"os"

	"github.com/flacatus/che-inspector/pkg/common/clog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/flacatus/che-inspector/pkg/api"
)

var env []string

func NewDockerTestSuite() {

}

// Comment
func RunTestsInDockerContainer(dockerClient *client.Client, testSpec *api.CheTestsSpec) (err error) {
	clog.LOGGER.Info("Pulling Test image...")
	if err := PullTestImage(dockerClient, testSpec); err != nil {
		return err
	}

	if err := CreateAndStartContainer(dockerClient, testSpec); err != nil {
		return err
	}

	return nil
}

// Comment
func PullTestImage(dockerClient *client.Client, testSpec *api.CheTestsSpec) (err error) {
	_, err = dockerClient.ImagePull(context.Background(), testSpec.Image, types.ImagePullOptions{})
	if err != nil {
		clog.LOGGER.Fatalf("Error pulling image %s, %v", testSpec.Image, err)

		return err
	}

	return nil
}

// Comment
func CreateAndStartContainer(dockerClient *client.Client, testSpec *api.CheTestsSpec) (err error) {
	for _, e := range testSpec.Env {
		env = append(env, e.Name+"="+e.Value)
	}
	if os.MkdirAll(testSpec.Artifacts.To, 0755) != nil && !os.IsExist(err) {
		clog.LOGGER.Fatalf("Cannot create directory %v", testSpec.Artifacts.To)
	}

	resp, err := dockerClient.ContainerCreate(context.Background(), &container.Config{
		Image: testSpec.Image,
		Env:   env,
		Tty:   false,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
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

// ExistDirectoryPath returns whether the given file or directory exists
func ExistDirectoryPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
