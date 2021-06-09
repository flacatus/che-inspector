package docker

import (
	"github.com/flacatus/che-inspector/pkg/api"
)

// Comment
func RunDockerSuite(instance *api.CliContext) (err error) {
	for _, suite := range instance.CheInspector.Spec.Tests {
		if err := RunTestsInDockerContainer(instance.DockerClient, &suite); err != nil {
			return err

		}
	}
	return nil
}
