package suites

import (
	"github.com/flacatus/che-inspector/pkg/common/instance"
	"github.com/flacatus/che-inspector/pkg/suites/docker"
	"github.com/flacatus/che-inspector/pkg/suites/kubernetes"
)

func RunTestSuite(instance *instance.CliContext) (err error) {
	for _, suite := range instance.CheInspector.Spec.Tests {
		if suite.ContainerContext == "docker" {
			return docker.RunDockerSuite(instance)
		}
		if suite.ContainerContext == "kubernetes" {
			return kubernetes.StartK8STestSuites(instance)
		}
	}
	return nil
}
