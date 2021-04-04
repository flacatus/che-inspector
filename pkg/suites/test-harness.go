package suites

import (
	"errors"
	"fmt"
	"time"

	"github.com/flacatus/che-inspector/pkg/common/client"
	"github.com/flacatus/che-inspector/pkg/common/instance"
	"github.com/flacatus/che-inspector/pkg/common/util"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployTestHarness(k8sClient *client.K8sClient, testSpec *instance.CheTestsSpec) (err error) {
	role, err := k8sClient.Kube().RbacV1().Roles(testSpec.Namespace).Create(getSpecRole(testSpec))
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof("Successufully create roles for test-harness %s.", role.Name)

	rb, err := k8sClient.Kube().RbacV1().RoleBindings(testSpec.Namespace).Create(getRoleBindingSpec(testSpec))
	if err != nil {
		return err
	}
	logrus.Infof("Successufully create roleBinding for test-harness %s.", rb.Name)

	pod, err := k8sClient.Kube().CoreV1().Pods(testSpec.Namespace).Create(GetTestHarnessPodSpec(testSpec))
	if err != nil {
		return err
	}
	terminated, err := waitForContainerToBeTerminated(k8sClient, testSpec, pod.Name)
	if terminated {
		err = util.CopyArtifactsFromPod(testSpec.Artifacts.FromContainerPath, testSpec.Artifacts.To, pod.Name, testSpec.Namespace, "download")
	} else {
		return fmt.Errorf("Failed to get test-harness pod status")
	}
	return err
}

func GetTestHarnessPodSpec(testSpec *instance.CheTestsSpec) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: testSpec.Name,
			Namespace:    testSpec.Namespace,
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: "test-run-results",
				},
			},
			RestartPolicy: "Never",
			Containers: []corev1.Container{
				{
					Name:            testSpec.Name,
					Image:           "quay.io/crw/osd-e2e:nightly",
					Args:            testSpec.Args,
					ImagePullPolicy: "Always",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "test-run-results",
							MountPath: "/test-run-results",
						},
					},
				},
				{
					Name:  "download",
					Image: "eeacms/rsync",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "test-run-results",
							MountPath: "/test-run-results",
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

func getSpecRole(testSpec *instance.CheTestsSpec) *v1.Role {
	return &v1.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Role",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-harness-role",
			Namespace: testSpec.Namespace,
		},
		Rules: []v1.PolicyRule{
			{
				APIGroups: []string{""},
				Verbs:     []string{"*"},
				Resources: []string{"pods", "services", "serviceaccounts", "endpoints", "persistentvolumeclaims", "events", "configmaps", "secrets", "pods/exec", "namespaces", "pods/log"},
			},
			{
				APIGroups: []string{"monitoring.coreos.com"},
				Verbs:     []string{"GET"},
				Resources: []string{"servicemonitors"},
			},
			{
				APIGroups: []string{"org.eclipse.che"},
				Verbs:     []string{"*"},
				Resources: []string{"checlusters", "checlusters/status", "checlusters/finalizers"},
			},
			{
				APIGroups: []string{"org.eclipse.che"},
				Verbs:     []string{"*"},
				Resources: []string{"checlusters", "checlusters/status", "checlusters/finalizers"},
			},
			{
				APIGroups: []string{"operators.coreos.com"},
				Verbs:     []string{"*"},
				Resources: []string{"subscriptions", "clusterserviceversions", "operatorgroups"},
			},
		},
	}
}

func getRoleBindingSpec(testSpec *instance.CheTestsSpec) *v1.RoleBinding {
	return &v1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-harness-role-binding",
			Namespace: testSpec.Namespace,
		},
		Subjects: []v1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "default",
				Namespace: testSpec.Namespace,
			},
		},
		RoleRef: v1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     "test-harness-role",
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
					return true, nil
				}
			}
		}
	}
}
