package cmd

import (
	"github.com/flacatus/che-inspector/pkg/cmd/version"
	"github.com/flacatus/che-inspector/pkg/util"
	"github.com/spf13/cobra"
	"io"
	cliflag "k8s.io/component-base/cli/flag"
	"os"
)

var pene bool

// NewCheInspectorCommand creates the `che-inspector` command with default arguments
func NewCheInspectorCommand() *cobra.Command {
	return NewCheInspectorCobraCommand(os.Stdin, os.Stdout, os.Stderr)
}

// NewCheInspectorCobraCommand creates the `che-inspector` command and its nested children.
func NewCheInspectorCobraCommand(in io.Reader, out, err io.Writer) *cobra.Command {
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
			},
		},
	}
	groups.Add(cmds)

	return cmds
}