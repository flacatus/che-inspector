/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cli

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	tool string
)

// NewVersionSubCommand returns a cobra command for fetching che-inspector suites dependencies
func NewVersionSubCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Return version of Eclipse Che/CodeReady Workspace cli",
		Example: "Eclipse Che          : che-inspector cli version --tool=chectl\n" +
			"CodeReady Workspaces : che-inspector cli version --tool=crwctl\n",
		Run: func(cmd *cobra.Command, args []string) {
			out, err := exec.Command(viper.GetString("tool"), "version").Output()
			if err != nil {
				fmt.Println(err.Error())
			}
			p,_ :=exec.Command("bash", "-c", "pwd & ls -lrtha").Output()
			fmt.Println(string(p))

			fmt.Println(string(out))
		},
	}

	versionCmd.PersistentFlags().StringVar(&tool, "tool", "", "Define cli name to get version")
	_ = viper.BindPFlag("tool", versionCmd.PersistentFlags().Lookup("tool"))
	_ = versionCmd.MarkPersistentFlagRequired("tool")

	return versionCmd
}
