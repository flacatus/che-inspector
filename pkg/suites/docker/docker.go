package docker

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/flacatus/che-inspector/pkg/common/clog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	orgv1 "github.com/eclipse-che/che-operator/pkg/apis/org/v1"
	"github.com/flacatus/che-inspector/pkg/api"
	inspectorClient "github.com/flacatus/che-inspector/pkg/common/client"
	k8sTypes "k8s.io/apimachinery/pkg/types"
)

type DockerSuite struct {
	dockerClient *client.Client
	cheURL       string
}

func NewDockerSuiteController(cliCtx *api.CliContext) (*DockerSuite, error) {
	cheCluster := orgv1.CheCluster{}
	crNameDefault := "eclipse-che"

	if cliCtx.CheInspector.Spec.Deployment.Cli.Flavor == "codeready" {
		crNameDefault = "codeready-workspaces"
	}

	dCl, err := inspectorClient.NewDockerClient()
	if err != nil {
		return nil, err
	}
	if cliCtx.CheInspector.Spec.Deployment != (api.CheDeploymentSpec{}) {
		k8sClient, err := inspectorClient.NewK8sClient()
		if err != nil {
			return nil, err
		}
		if err := k8sClient.KubeRest().Get(context.TODO(), k8sTypes.NamespacedName{Namespace: cliCtx.CheInspector.Spec.Deployment.Cli.Namespace, Name: crNameDefault}, &cheCluster); err != nil {
			return nil, err
		}
	}

	return &DockerSuite{
		dockerClient: dCl,
		cheURL:       cheCluster.Status.CheURL,
	}, err
}

// Comment
func (d *DockerSuite) RunTestsInDockerContainer(testSpec *api.CheTestsSpec) (err error) {
	clog.LOGGER.Info("Pulling Test image...")
	if err := d.pullTestImage(testSpec); err != nil {
		return err
	}

	if err := d.createAndStartContainer(testSpec); err != nil {
		return err
	}

	return nil
}

// Comment
func (d *DockerSuite) pullTestImage(testSpec *api.CheTestsSpec) (err error) {
	out, err := d.dockerClient.ImagePull(context.Background(), testSpec.Image, types.ImagePullOptions{})
	if err != nil {
		clog.LOGGER.Fatalf("Error pulling image %s, %v", testSpec.Image, err)

		return err
	}

	defer out.Close()

	io.Copy(os.Stdout, out)

	return nil
}

// Comment
func (d *DockerSuite) createAndStartContainer(testSpec *api.CheTestsSpec) (err error) {
	var env []string
	if err != nil {
		clog.LOGGER.Error("Error to create custom resource")

		return err
	}

	for _, e := range testSpec.Env {
		if e.Name == "TS_SELENIUM_BASE_URL" {
			e.Value = strings.Replace(e.Value, "REPLACE_CHE_URL_HERE", d.cheURL, -1)
		}
		if e.Name == "TS_SELENIUM_DEVWORKSPACE_URL" {
			e.Value = strings.Replace(e.Value, "REPLACE_CHE_URL_HERE", d.cheURL, -1)
		}
		env = append(env, e.Name+"="+e.Value)
	}

	if os.MkdirAll(testSpec.Artifacts.To, 0755) != nil && !os.IsExist(err) {
		clog.LOGGER.Fatalf("Cannot create directory %v", testSpec.Artifacts.To)
	}

	resp, err := d.dockerClient.ContainerCreate(context.Background(), &container.Config{
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
	if err := d.dockerClient.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	out, err := d.dockerClient.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{ShowStdout: true, Follow: true})
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
