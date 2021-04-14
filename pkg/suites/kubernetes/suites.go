package kubernetes

import (
	"errors"
	"fmt"
	"github.com/flacatus/che-inspector/pkg/common/client"
	"github.com/flacatus/che-inspector/pkg/common/clog"
	"github.com/flacatus/che-inspector/pkg/common/instance"
	"github.com/flacatus/che-inspector/pkg/common/reporter"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const (
	artifactsVolumeName                 = "test-run-results"
	artifactsDownloadContainerName      = "download"
	downloadArtifactsContainerImage     = "eeacms/rsync"
	testHarnessRoleName                 = "test-harness-role"
	testHarnessRoleBindingName          = "test-harness-role-binding"
)

func StartK8STestSuites(instance *instance.CliContext) (err error) {
	// TODO: Refactor this function
	for _, suite := range instance.CheInspector.Spec.Tests {
		if suite.Name == "test-harness" {
			clog.LOGGER.Info("Starting test-harness test suite")

			return DeployTestHarness(instance.Client, &suite)
		}
		if suite.Name == "happy-path" {
			clog.LOGGER.Info("Starting happy path test suite")

			if err := DeployHappyPath(instance.Client, &suite); err != nil {
				_ = reporter.SendSlackMessage(&instance.CheInspector.Spec.Report[0], "Che-Inspector: Failed to run hapy path tests ")

				return err
			}
			_ = reporter.SendSlackMessage(&instance.CheInspector.Spec.Report[0], "Che-Inspector: DevWorkspace run successfully in openshift ")
		}
	}

	return nil
}

// The GetTestSuitePodSpec returns pod specification with pod to be created with go client
func GetTestSuitePodSpec(testSpec *instance.CheTestsSpec) *corev1.Pod  {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: testSpec.Name,
			Namespace: testSpec.Namespace,
		},
		Spec:       corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: artifactsVolumeName,
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name: testSpec.Name,
					Image: testSpec.Image,
					Args: testSpec.Args,
					Env: testSpec.Env,
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

func waitForContainerToBeTerminated(k8sClient *client.K8sClient, testSpec *instance.CheTestsSpec, podName string) (terminated bool, err error) {
	for {
		select {
		case <-time.After(15 * time.Minute):
			return false, errors.New("timed out")
		case <-time.Tick(15 * time.Second):
			pod, err := k8sClient.Kube().CoreV1().Pods(testSpec.Namespace).Get(podName, metav1.GetOptions{})
			if err != nil {
				return true, err
			}
			for _, container := range pod.Status.ContainerStatuses {
				if container.Name == testSpec.Name && container.State.Terminated != nil {
					fmt.Println(container.State.Running.StartedAt)
					return true, nil
				}
			}
		}
	}
}
