package docker

import (
	"github.com/flacatus/che-inspector/pkg/api"
	"github.com/flacatus/che-inspector/pkg/common/reporter"
)

func RunDockerSuite(instance *api.CliContext) (err error) {
	for _, suite := range instance.CheInspector.Spec.Tests {
		if err := RunHappyPathDocker(instance.DockerClient, &suite); err != nil {
			_ = reporter.SendSlackMessage(&instance.CheInspector.Spec.Report[0], "Che-Inspector: DevWorkspace run failed in openshift")
			return err

		}
		_ = reporter.SendSlackMessage(&instance.CheInspector.Spec.Report[0], "Che-Inspector: DevWorkspace run successfully in openshift")
	}
	return nil
}
