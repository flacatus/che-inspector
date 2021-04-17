package api

import (
	"fmt"
	"io/ioutil"

	dockerClient "github.com/docker/docker/client"
	"github.com/flacatus/che-inspector/pkg/common/client"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Comment
type CliContext struct {
	CheInspector *CheInspector
	Client       *client.K8sClient
	DockerClient *dockerClient.Client
}

// Comment
func GetCliContext() (cliContext *CliContext, e error) {
	cliInstance, err := readCheInspectorConfig()
	if err != nil {
		return nil, err
	}

	k8sClient, err := client.NewK8sClient()
	if err != nil {
		return nil, err
	}

	dockerCl, err := client.NewDockerClient()
	if err != nil {
		return nil, err
	}
	return &CliContext{
		CheInspector: cliInstance,
		Client:       k8sClient,
		DockerClient: dockerCl,
	}, nil
}

// Comment
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
