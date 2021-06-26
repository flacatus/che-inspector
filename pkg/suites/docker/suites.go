package docker

import (
	"github.com/flacatus/che-inspector/pkg/api"
	"github.com/flacatus/che-inspector/pkg/common/client"
)

// Comment
func RunDockerSuite(instance *api.CliContext) (err error) {
	instance.DockerClient, err = client.NewDockerClient()
	if err != nil {
		return err
	}

	for _, suite := range instance.CheInspector.Spec.Tests {
		if err := RunTestsInDockerContainer(instance.DockerClient, &suite); err != nil {
			return err
		}
	}

	return nil
}
