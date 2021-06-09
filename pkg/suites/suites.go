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

package suites

import (
	"github.com/flacatus/che-inspector/pkg/api"
	"github.com/flacatus/che-inspector/pkg/suites/docker"
	"github.com/flacatus/che-inspector/pkg/suites/kubernetes"
)

// Comment
func RunTestSuite(instance *api.CliContext) (err error) {
	for _, suite := range instance.CheInspector.Spec.Tests {
		if suite.ContainerContext == "docker" {
			return docker.RunDockerSuite(instance)
		}
		if suite.ContainerContext == "kubernetes" {
			return kubernetes.StartK8STestSuites(instance)
		}
	}
	return nil
}
