// Copyright (c) 2021 The Jaeger Authors.
// //
// // Copyright (c) 2021 Red Hat, Inc.
// // This program and the accompanying materials are made
// // available under the terms of the Eclipse Public License 2.0
// // which is available at https://www.eclipse.org/legal/epl-2.0/
// //
// // SPDX-License-Identifier: EPL-2.0
// //
// // Contributors:
// //   Red Hat, Inc. - initial API and implementation
// //

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
				run.NewRunCommand(),
				reporter.NewReporterCommand(),
				version.AddVersionCommand(),
			},
		},
	}
	groups.Add(cmds)

	return cmds
}
