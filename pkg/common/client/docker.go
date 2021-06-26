package client

import (
	"github.com/docker/docker/client"
)

// NewDockerClient creates docker client wrapper with helper functions to talk docker API
func NewDockerClient() (*client.Client, error) {
	clientCfg, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return clientCfg, nil
}
