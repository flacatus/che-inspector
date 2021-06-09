package run

import (
	"github.com/flacatus/che-inspector/pkg/api"
	"github.com/flacatus/che-inspector/pkg/common/clog"
	"github.com/flacatus/che-inspector/pkg/common/validator"
	"github.com/flacatus/che-inspector/pkg/suites"
	_ "github.com/flacatus/che-inspector/pkg/suites"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var file string

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start test suites from a custom configuration",
		Long: `
      Start test suites from a given Che Inspector file.

      Find more information at:
            In PROGRESS`,
		Example: "che-inspector run --file=samples/happy-path.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			context, err := PreRunningTasks()
			if err != nil {
				clog.LOGGER.Fatal(err)
			}
			_ = suites.RunTestSuite(context)
			if err != nil {
				clog.LOGGER.Fatal(err)
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Configuration file with definitions of all suites cases to run against Che")
	_ = viper.BindPFlag("file", cmd.PersistentFlags().Lookup("file"))
	_ = cmd.MarkPersistentFlagRequired("file")

	return cmd
}

func PreRunningTasks() (cliContext *api.CliContext, err error) {
	clog.LOGGER.Infof("Starting to validate test configuration and generating cli context")
	context, err := api.GetCliContext()
	if err != nil {
		return context, err
	}

	err = validator.CheInspectorValidator(context.CheInspector)
	if err != nil {
		return context, err
	}

	return context, nil
}
