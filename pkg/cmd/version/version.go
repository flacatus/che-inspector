package version

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"runtime"
)

var (
	// cliVersion is the constant representing the version of the che-inspector binary
	cliVersion = "unknown"
	// gitCommit is a constant representing the source version that
	// generated this build. It should be set during build via -ldflags.
	gitCommit string
	// buildDate in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	buildDate string
)

type Version struct {
	CliVersion string `yaml:"cliVersion"`
	GitCommit  string `yaml:"gitCommit"`
	BuildDate  string `yaml:"buildDate"`
	GoOs       string `yaml:"goOs"`
	GoArch     string `yaml:"goArch"`
}

func getVersion() Version {
	return Version{
		CliVersion: cliVersion,
		GitCommit:  gitCommit,
		BuildDate:  buildDate,
		GoOs:       runtime.GOOS,
		GoArch:     runtime.GOARCH,
	}
}

func (v Version) Print() {
	version, _ := yaml.Marshal(v)
	fmt.Printf(string(version) + "\n")
}

func AddVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Print the che-inspector version",
		Long:    `
      Check che-inspector version.`,
		Example: `che-inspector version`,
		Run:     runVersion,
	}
}

func runVersion(_ *cobra.Command, _ []string) {
	getVersion().Print()
}
