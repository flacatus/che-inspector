package reporter

import (
	"fmt"

	"github.com/flacatus/che-inspector/pkg/api"
	report_portal "github.com/flacatus/che-inspector/pkg/api/report-portal"
	"github.com/flacatus/che-inspector/pkg/common/clog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var file string

func NewReporterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Send test report to report-portal",
		Long: `
		Send test report to report-portal.

      Find more information at:
            In PROGRESS`,
		Example: "che-inspector report --file=samples/happy-path.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			context, err := api.GetCliContext()
			if err != nil {
				clog.LOGGER.Info("Tests failed")
			}
			fmt.Println(context.CheInspector)
			for _, r := range context.CheInspector.Spec.Report {
				if r.ReportPortal.Name != "" {
					clientReport := report_portal.NewReportPortalClient(&r.ReportPortal)
					clientReport.SendResultsToReportPortal()
				}
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Configuration file with definitions of all suites cases to run against Che")
	_ = viper.BindPFlag("file", cmd.PersistentFlags().Lookup("file"))
	_ = cmd.MarkPersistentFlagRequired("file")

	return cmd
}

type Reporter struct{}

func NewReport() *Reporter {
	return &Reporter{}
}
