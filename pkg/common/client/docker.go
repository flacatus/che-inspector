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
