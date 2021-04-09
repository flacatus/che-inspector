package suites

import (
	"fmt"
	"github.com/flacatus/che-inspector/pkg/common/client"
	"github.com/flacatus/che-inspector/pkg/common/instance"
	"github.com/flacatus/che-inspector/pkg/common/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployHappyPath(k8sClient *client.K8sClient, testSpec *instance.CheTestsSpec) (err error) {
	pod, err := k8sClient.Kube().CoreV1().Pods(testSpec.Namespace).Create(GetHappyPathPodSpec(testSpec))
	if err != nil {
		return err
	}

	terminated, err := waitForContainerToBeTerminated(k8sClient, testSpec, pod.Name)
	if terminated {
		err = util.CopyArtifactsFromPod(testSpec.Artifacts.FromContainerPath, testSpec.Artifacts.To, pod.Name, testSpec.Namespace, artifactsDownloadContainerName)
	} else {
		return fmt.Errorf("Failed to get test-harness pod status")
	}

	return err
}

func GetHappyPathPodSpec(testSpec *instance.CheTestsSpec) *corev1.Pod {
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
			RestartPolicy: "Never",
			Containers: []corev1.Container{
				{
					Name:  testSpec.Name,
					Image: testSpec.Image,
					Env:   testSpec.Env,
					ImagePullPolicy: "Always",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      artifactsVolumeName,
							MountPath: happyPathArtifactsVolumeMonthPath,
						},
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("2"),
							corev1.ResourceMemory: resource.MustParse("3Gi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("2"),
							corev1.ResourceMemory: resource.MustParse("4Gi"),
						},
					},
				},
				{
					Name:  artifactsDownloadContainerName,
					Image: downloadArtifactsContainerImage,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      artifactsVolumeName,
							MountPath: happyPathArtifactsVolumeMonthPath,
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
