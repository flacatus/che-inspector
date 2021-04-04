package suites

import (
	"github.com/flacatus/che-inspector/pkg/common/instance"
	"github.com/sirupsen/logrus"
)

func StartCheSuiteTests(instance *instance.CliContext) error {
	for _, suite := range instance.CheInspector.Spec.Tests {
		if suite.Name == "test-harness" {
			logrus.Info("Starting test-harness test suite")
			return DeployTestHarness(instance.Client, &suite)
		}
	}
	return nil
}
