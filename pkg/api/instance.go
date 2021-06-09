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

	return &CliContext{
		CheInspector: cliInstance,
		Client:       nil,
		DockerClient: nil,
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
