package che

import (
	"path/filepath"

	"github.com/flacatus/che-inspector/pkg/common/clog"
	"github.com/flacatus/che-inspector/pkg/util"
)

// DeployIde: method which deploy Che/CRW depending on which Che flavour is defined in CheInspector configurations.
// Possible flavour options are: codeready to deploy Red Hat CodeReady Workspaces using crwctl cli: https://github.com/redhat-developer/codeready-workspaces-chectl
// And other flavour option is che which deploy Eclipse Che using chectl cli: https://github.com/che-incubator/chectl
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

// Execute deployment command from a given of the ide cli and a given deployment flags
// Examples of Che/CRW deploy is: <cli-path> server:deploy --platform=openshift
func executeCheDeployCommand(cliPath string, deployFlags string) error {
	command := cliPath + " " + deployFlags
	return util.ExecuteBashCommand(command)
}
