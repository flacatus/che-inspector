package reporter

import (
	"github.com/flacatus/che-inspector/pkg/api"
	report_portal "github.com/flacatus/che-inspector/pkg/api/report-portal"
	"github.com/flacatus/che-inspector/pkg/common/clog"
	"github.com/spf13/cobra"
)

func NewReporterCommand() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "report",
		Short: "Send test report to report-portal",
		Long: `
		[WORK IN PROGRESS]Send test report to report-portal.

      Find more information at:
            In PROGRESS`,
		Example: "che-inspector report --file=samples/happy-path.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			context, err := api.GetCliContext(configFile)
			if err != nil {
				clog.LOGGER.Info(err)
			}
			for _, r := range context.CheInspector.Spec.Report {
				if r.ReportPortal.Name != "" {
					clientReport := report_portal.NewReportPortalClient(&r.ReportPortal)
					clientReport.SendResultsToReportPortal()
				}
			}
		},
	}

	cmd.PersistentFlags().StringVarP(&configFile, "file", "f", "", "Configuration file with definitions of all suites cases to run against Che")
	_ = cmd.MarkPersistentFlagRequired("file")

	return cmd
}
