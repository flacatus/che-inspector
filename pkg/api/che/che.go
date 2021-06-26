package che

import (
	"github.com/flacatus/che-inspector/pkg/api"
)

const (
	CHECTL_COMMAND_NAME    = "chectl"
	CHE_FLAVOUR_NAME       = "che"
	CODEREADY_FLAVOUR_NAME = "codeready"
	CHECTL_INSTALL_COMMAND = "bash <(curl -sL  https://www.eclipse.org/che/chectl/) --channel=next"
)

type DeployChe struct {
	deployAPI api.CheDeploymentSpec
}

func NewCheController(deploySpec api.CheDeploymentSpec) *DeployChe {
	return &DeployChe{
		deployAPI: deploySpec,
	}
}
