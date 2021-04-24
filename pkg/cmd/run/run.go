package run

import (
	"github.com/flacatus/che-inspector/pkg/api"
	report_portal "github.com/flacatus/che-inspector/pkg/api/report-portal"
	"github.com/flacatus/che-inspector/pkg/common/clog"
	"github.com/flacatus/che-inspector/pkg/common/validator"
	"github.com/flacatus/che-inspector/pkg/suites"
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

			context, err := api.GetCliContext()

			rp := &context.CheInspector.Spec.Report[0].ReportPortal
			clientReport := report_portal.NewReportPortalClient(rp)
			clientReport.SendResultsToReportPortal()

			if err != nil {
				clog.LOGGER.Fatal(err)
			}

			err = validator.CheInspectorValidator(context.CheInspector)
			if err != nil {
				clog.LOGGER.Fatal(err)
				return
			}
			clog.LOGGER.Infof("Successfully validated configuration")
			err = suites.RunTestSuite(context)
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
