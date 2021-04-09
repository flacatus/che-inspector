package suites

import (
	"github.com/flacatus/che-inspector/pkg/common/clog"
	"github.com/flacatus/che-inspector/pkg/common/instance"
	"github.com/flacatus/che-inspector/pkg/common/reporter"
)

const (
	artifactsVolumeName                 = "test-run-results"
	testHarnessArtifactsVolumeMonthPath = "/test-run-results"
	artifactsDownloadContainerName      = "download"
	downloadArtifactsContainerImage     = "eeacms/rsync"
	testHarnessRoleName                 = "test-harness-role"
	testHarnessRoleBindingName          = "test-harness-role-binding"
	happyPathArtifactsVolumeMonthPath   = "/tmp/e2e/report/"
)

func StartCheSuiteTests(instance *instance.CliContext) (err error) {
	// TODO: Refactor this function
	for _, suite := range instance.CheInspector.Spec.Tests {
		if suite.Name == "test-harness" {
			clog.LOGGER.Info("Starting test-harness test suite")

			return DeployTestHarness(instance.Client, &suite)
		}
		if suite.Name == "happy-path" {
			clog.LOGGER.Info("Starting happy path test suite")

			if err := DeployHappyPath(instance.Client, &suite); err != nil {
				_ = reporter.SendSlackMessage(&instance.CheInspector.Spec.Report[0], "Che-Inspector: Failed to run hapy path tests ")

				return err
			}
			_ = reporter.SendSlackMessage(&instance.CheInspector.Spec.Report[0], "Che-Inspector: DevWorkspace run successfully in openshift ")
		}
	}

	return nil
}
