package che

import (
	"path/filepath"

	"github.com/flacatus/che-inspector/pkg/common/clog"
	"github.com/flacatus/che-inspector/pkg/util"
)

func (d *DeployChe) DeployIde() error {
	if d.deployAPI.Cli.Flavor == CHE_FLAVOUR_NAME {
		clog.LOGGER.Infof("Start to deploy Eclipse Che: %s", CHECTL_COMMAND_NAME+" "+d.deployAPI.Cli.Flags)
		if err := executeCheDeployCommand(CHECTL_COMMAND_NAME, d.deployAPI.Cli.Flags); err != nil {
			return err
		}
	}

	if d.deployAPI.Cli.Flavor == CODEREADY_FLAVOUR_NAME {
		ansoluteRoute, _ := filepath.Abs(d.deployAPI.Cli.InstallPath)
		crwctlCompletePath := ansoluteRoute + "/crwctl/bin/crwctl"

		clog.LOGGER.Infof("Start to deploy Red Hat CodeReady Workspaces: %s", crwctlCompletePath+" "+d.deployAPI.Cli.Flags)
		if err := executeCheDeployCommand(crwctlCompletePath, d.deployAPI.Cli.Flags); err != nil {
			return err
		}
	}
	return nil
}

func executeCheDeployCommand(crwctlPath string, deployFlags string) error {
	command := crwctlPath + " " + deployFlags

	return util.ExecuteBashCommand(command)
}
