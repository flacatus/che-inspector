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

package kubernetes

import (
	"context"
	e "errors"
	"fmt"
	"time"

	"github.com/flacatus/che-inspector/pkg/api"
	"github.com/flacatus/che-inspector/pkg/common/client"
	"github.com/flacatus/che-inspector/pkg/common/clog"
	"github.com/flacatus/che-inspector/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	artifactsVolumeName             = "test-run-results"
	artifactsDownloadContainerName  = "download"
	downloadArtifactsContainerImage = "eeacms/rsync"
	testHarnessRoleName             = "test-harness-role"
	testHarnessRoleBindingName      = "test-harness-role-binding"
)

// Comment
func StartK8STestSuites(instance *api.CliContext) (err error) {
	for _, suite := range instance.CheInspector.Spec.Tests {
		if err := DeployTestSuite(instance.Client, &suite); err != nil {

			return err
		}
	}

	return nil
}

// Comment
func DeployTestSuite(k8sClient *client.K8sClient, testSpec *api.CheTestsSpec) (err error) {
	if _, err := k8sClient.Kube().CoreV1().Namespaces().Get(context.TODO(), testSpec.Namespace, metav1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			clog.LOGGER.Infof("Namespace %s doesn't exist. Creating new one...", testSpec.Namespace)
			if _, err := k8sClient.Kube().CoreV1().Namespaces().Create(context.TODO(), GetNamespaceSpec(testSpec), metav1.CreateOptions{}); err != nil {
				clog.LOGGER.Fatalf("Failed to create namespace %s: %v", testSpec.Namespace, err)
			}
		}
	}

	pod, err := k8sClient.Kube().CoreV1().Pods(testSpec.Namespace).Create(context.TODO(), GetTestSuitePodSpec(testSpec), metav1.CreateOptions{})
	if err != nil {
		return err
	}

	terminated, _ := waitForContainerToBeTerminated(k8sClient, testSpec, pod.Name)
	if terminated {
		util.CopyArtifactsFromPod(testSpec.Artifacts.FromContainerPath, testSpec.Artifacts.To, pod.Name, testSpec.Namespace, artifactsDownloadContainerName)
	} else {
		return fmt.Errorf("Failed to get test pod status")
	}

	return util.ExecInContainer(pod.Name, artifactsDownloadContainerName, testSpec.Namespace, "touch /tmp/done")
}

// GetNamespaceSpec return namespace object
func GetNamespaceSpec(testSpec *api.CheTestsSpec) *corev1.Namespace {
	return &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: testSpec.Namespace,
		},
	}
}

// The GetTestSuitePodSpec returns pod specification with pod to be created with go client
func GetTestSuitePodSpec(testSpec *api.CheTestsSpec) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: testSpec.Name,
			Namespace:    testSpec.Namespace,
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: artifactsVolumeName,
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:            testSpec.Name,
					Image:           testSpec.Image,
					Args:            testSpec.Args,
					Env:             testSpec.Env,
					ImagePullPolicy: corev1.PullAlways,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      artifactsVolumeName,
							MountPath: testSpec.Artifacts.FromContainerPath,
						},
					},
				},
				{
					Name:  artifactsDownloadContainerName,
					Image: downloadArtifactsContainerImage,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      artifactsVolumeName,
							MountPath: testSpec.Artifacts.FromContainerPath,
						},
					},
					Command: []string{"sh"},
					Args: []string{
						"-c",
						"while true; if [[ -f /tmp/done ]]; then exit 0; fi; do sleep 1; done",
					},
				},
			},
		},
	}
}

// TODO: Use wait.poll
func waitForContainerToBeTerminated(k8sClient *client.K8sClient, testSpec *api.CheTestsSpec, podName string) (terminated bool, err error) {
	for {
		select {
		case <-time.After(15 * time.Minute):
			return false, e.New("timed out")
		case <-time.Tick(15 * time.Second):
			pod, err := k8sClient.Kube().CoreV1().Pods(testSpec.Namespace).Get(context.TODO(), podName, metav1.GetOptions{})
			if err != nil {
				return true, err
			}
			for _, container := range pod.Status.ContainerStatuses {
				if container.Name == testSpec.Name && container.State.Terminated != nil {
					return true, nil
				}
			}
		}
	}
}
