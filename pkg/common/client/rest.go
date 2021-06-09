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
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Comment
func InitRestClient() (*rest.Config, error, *corev1client.CoreV1Client) {
	// Instantiate loader for kubeconfig file.
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Get a rest.Config from the kubeConfig file.  This will be passed into all
	// the client objects we create.
	restConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	// Create a Kubernetes core/v1 client.
	coreClient, err := corev1client.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	return restConfig, err, coreClient
}
