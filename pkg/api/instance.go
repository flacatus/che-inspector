package api

import (
	"fmt"
	"io/ioutil"

	dockerClient "github.com/docker/docker/client"
	"github.com/flacatus/che-inspector/pkg/common/client"
	"gopkg.in/yaml.v2"
)

// Comment
type CliContext struct {
	CheInspector *CheInspector
	Client       *client.K8sClient
	DockerClient *dockerClient.Client
}

// Comment
func GetCliContext(configFile string) (cliContext *CliContext, e error) {
	cliInstance, err := readCheInspectorConfig(configFile)
	if err != nil {
		return nil, err
	}

	return &CliContext{
		CheInspector: cliInstance,
		Client:       nil,
		DockerClient: nil,
	}, nil
}

// Comment
func readCheInspectorConfig(configFile string) (instance *CheInspector, e error) {
	c := &CheInspector{}
	yamlFileContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Error on getting configuration file: %v", err)
	}
	err = yaml.Unmarshal(yamlFileContent, c)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling configuration file : %v", err)
	}

	return c, nil
}
