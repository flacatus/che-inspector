package run

import (
	"github.com/flacatus/che-inspector/pkg/api"
	"github.com/flacatus/che-inspector/pkg/api/che"
	"github.com/flacatus/che-inspector/pkg/common/clog"
	"github.com/flacatus/che-inspector/pkg/common/validator"
	"github.com/flacatus/che-inspector/pkg/suites"
	_ "github.com/flacatus/che-inspector/pkg/suites"
	"github.com/spf13/cobra"
)

func NewRunCommand() *cobra.Command {
	var fileStream string
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Send test report to report-portal",
		Long: `
		Send test report to report-portal.

      Find more information at:
            In PROGRESS`,
		Example: "che-inspector report --file=samples/happy-path.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			context, err := PreRunningTasks(fileStream)
			if err != nil {
				clog.LOGGER.Info(err)
			}
			if err := suites.RunTestSuite(context); err != nil {
				clog.LOGGER.Fatal(err)
			}
		},
	}

	cmd.PersistentFlags().StringVarP(&fileStream, "file", "f", "", "Configuration file with definitions of all suites cases to run against Che")
	_ = cmd.MarkPersistentFlagRequired("file")

	return cmd
}

func PreRunningTasks(configFile string) (cliContext *api.CliContext, err error) {
	clog.LOGGER.Infof("Starting to validate test configuration and generating cli context")
	context, err := api.GetCliContext(configFile)
	if err != nil {
		return context, err
	}

	err = validator.CheInspectorValidator(context.CheInspector)
	if err != nil {
		clog.LOGGER.Fatal(err)
	}

	if context.CheInspector.Spec.Deployment != (api.CheDeploymentSpec{}) {
		reconcileInstall := che.NewCheController(context.CheInspector.Spec.Deployment)
		if err = reconcileInstall.InstallCheCli(); err != nil {
			clog.LOGGER.Fatal(err)
		}

		if err = reconcileInstall.DeployIde(); err != nil {
			clog.LOGGER.Fatalf("Unable to deploy IDE. Please check your deploy command and flags, %s", err)
		}
	}

	return context, nil
}
