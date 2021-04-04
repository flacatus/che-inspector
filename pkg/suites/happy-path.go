package suites

import (
	"fmt"
	"github.com/flacatus/che-inspector/pkg/common/util"

	"github.com/flacatus/che-inspector/pkg/common/instance"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployHappyPath(cliContext *instance.CliContext,happyPathImage string, env []corev1.EnvVar) {
	env = append(env, corev1.EnvVar{
		Name: "BASE_URL",
		Value: "https://google.com"})

	pod, _ := cliContext.Client.Kube().CoreV1().Pods("happy-path").Create(GetHappyPathPodSpec(happyPathImage, env))
	err := util.CopyArtifactsFromPod("/dev", "/home/flacatusu/WORKSPACE/che-inspector/fla","che-df5d95985-767wk", "eclipse-che", "che")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pod.Name)
}

func GetHappyPathPodSpec(happyPathImage string, env []corev1.EnvVar) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "happy-path",
			Namespace:    "happy-path",
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: "test-run-results",
				},
			},
			Containers: []corev1.Container{
				{
					Name:  "happy-path",
					Image: happyPathImage,
					Env:   env,
					ImagePullPolicy: "Always",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "test-run-results",
							MountPath: "/tmp/e2e/report/",
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
			},
		},
	}
}
