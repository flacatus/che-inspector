package suites

import (
	"fmt"

	"github.com/flacatus/che-inspector/pkg/api"
	"github.com/flacatus/che-inspector/pkg/suites/docker"
)

// Comment
func RunTestSuite(instance *api.CliContext) (err error) {
	for _, suite := range instance.CheInspector.Spec.Tests {
		if suite.ContainerContext == "docker" {

			return RunDockerSuite(instance)
		}
		if suite.ContainerContext == "kubernetes" {
			return fmt.Errorf("Run tests in kubernetes not supported yet on che-inspector. Soon will be available")
		}
	}
	return nil
}

func RunDockerSuite(instance *api.CliContext) (err error) {
	for _, suite := range instance.CheInspector.Spec.Tests {
		dockerSuite, err := docker.NewDockerSuiteController(instance)
		if err != nil {
			return err
		}
		if err := dockerSuite.RunTestsInDockerContainer(&suite); err != nil {
			return err
		}
	}

	return nil
}
