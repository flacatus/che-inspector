package instance

import (
	"fmt"
	"io/ioutil"

	"github.com/flacatus/che-inspector/pkg/common/client"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type CliContext struct {
	CheInspector *CheInspector
	Client       *client.K8sClient
}

func GetCliContext() (cliContext *CliContext, e error) {
	cliInstance, err := readCheInspectorConfig()
	if err != nil {
		return nil, err
	}

	k8sClient, err := client.NewK8sClient()
	if err != nil {
		return nil, err
	}

	return &CliContext{
		CheInspector: cliInstance,
		Client:       k8sClient,
	}, nil
}

func readCheInspectorConfig() (instance *CheInspector, e error) {
	c := &CheInspector{}
	yamlFileContent, err := ioutil.ReadFile(viper.GetString("file"))
	if err != nil {
		return nil, fmt.Errorf("Error on getting configuration file: %v", err)
	}
	err = yaml.Unmarshal(yamlFileContent, c)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling configuration file : %v", err)
	}

	return c, nil
}
