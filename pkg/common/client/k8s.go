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
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// Comment
type K8sClient struct {
	kubeClient *kubernetes.Clientset
}

// NewK8sClient creates kubernetes client wrapper with helper functions to talk with API Server
func NewK8sClient() (*K8sClient, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &K8sClient{
		kubeClient: client,
	}, nil
}

// Kube returns the clientset for Kubernetes upstream.
func (c *K8sClient) Kube() kubernetes.Interface {
	return c.kubeClient
}
