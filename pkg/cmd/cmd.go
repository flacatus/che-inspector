package cmd

import (
	"github.com/flacatus/che-inspector/pkg/cmd/reporter"
	"github.com/flacatus/che-inspector/pkg/cmd/run"
	"github.com/flacatus/che-inspector/pkg/cmd/version"
	"github.com/flacatus/che-inspector/pkg/util"
	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
)

// NewCheInspectorCobraCommand creates the `che-inspector` command and its nested children.
func NewCheInspectorCobraCommand() *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "che-inspector",
		Short: "Che Inspector deploy tests anywhere",
		Long: `
      Che Inspector controls test workflow management.

      Find more information at:
            In PROGRESS`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	flags := cmds.PersistentFlags()
	flags.SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc) // Warn for "_" flags

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	// From this point and forward we get warnings on flags that contain "_" separators
	cmds.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)

	groups := util.CommandGroups{
		{
			Message: "Check Che Inspector version",
			Commands: []*cobra.Command{
				version.AddVersionCommand(),
				run.NewRunCommand(),
				reporter.NewReporterCommand(),
			},
		},
	}
	groups.Add(cmds)

	return cmds
}
