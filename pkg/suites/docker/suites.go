package docker

import (
	"github.com/flacatus/che-inspector/pkg/common/client"
	"github.com/flacatus/che-inspector/pkg/common/instance"
)

func RunDockerSuite(instance *instance.CliContext) (err error){
	dockerClient, err := client.NewDockerClient()
	if  err != nil {
		return err
	}

	for _, suite := range instance.CheInspector.Spec.Tests {
		if err := RunHappyPathDocker(dockerClient, &suite); err != nil {
			return err
		}
	}
	return nil
}
