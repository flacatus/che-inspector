package che

import (
	"os"
	"path/filepath"

	"github.com/cavaliercoder/grab"
	"github.com/flacatus/che-inspector/pkg/common/clog"
	"github.com/flacatus/che-inspector/pkg/util"
)

// InstallCheCli install the cli which manage Red Hat CodeReady Workspaces or Eclipse Che depending on which
// flavour was specify in the config samples
func (d *DeployChe) InstallCheCli() error {
	if d.deployAPI.Cli.Flavor == CHE_FLAVOUR_NAME && !util.IsCommandAvailable(CHECTL_COMMAND_NAME) {
		if err := installChectl(); err != nil {
			return err
		}
	}

	if d.deployAPI.Cli.Flavor == CODEREADY_FLAVOUR_NAME {
		ansoluteRoute, _ := filepath.Abs(d.deployAPI.Cli.InstallPath)
		crwctlCompletePath := ansoluteRoute + "/crwctl/bin/crwctl"

		if !util.IsCommandAvailable(crwctlCompletePath) {
			return downloadAndUnzipCrwctl(d.deployAPI.Cli.InstallPath, d.deployAPI.Cli.Source)
		}
	}

	return nil
}

// Unzip crwctl code to a given directory
func downloadAndUnzipCrwctl(installPath string, source string) error {
	absoluteInstallPath, err := filepath.Abs(installPath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(absoluteInstallPath); os.IsNotExist(err) {
		if err = os.Mkdir(absoluteInstallPath, 0755); err != nil {
			return err
		}
	}

	resp, err := grab.Get(absoluteInstallPath, source)
	if err != nil {
		return err
	}
	if err = util.Untar(resp.Filename, absoluteInstallPath); err != nil {
		return err
	}

	return nil
}

// Install Eclipse Che using instructions from: https://github.com/che-incubator/chectl#installation
func installChectl() error {
	if util.IsCommandAvailable(CHECTL_COMMAND_NAME) {
		clog.LOGGER.Info("chectl is already installed")
		return nil
	}
	return util.ExecuteBashCommand(CHECTL_INSTALL_COMMAND)
}
