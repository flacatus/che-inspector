package kubernetes

import (
	"fmt"
	"github.com/flacatus/che-inspector/pkg/api"
	"github.com/flacatus/che-inspector/pkg/common/client"
	"github.com/flacatus/che-inspector/pkg/common/util"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployTestHarness(k8sClient *client.K8sClient, testSpec *api.CheTestsSpec) (err error) {
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

	pod, err := k8sClient.Kube().CoreV1().Pods(testSpec.Namespace).Create(GetTestSuitePodSpec(testSpec))
	if err != nil {
		return err
	}

	terminated, err := waitForContainerToBeTerminated(k8sClient, testSpec, pod.Name)
	if terminated {
		err = util.CopyArtifactsFromPod(testSpec.Artifacts.FromContainerPath, testSpec.Artifacts.To, pod.Name, testSpec.Namespace, artifactsDownloadContainerName)
	} else {
		return fmt.Errorf("Failed to get test-harness pod status   ")
	}

	return err
}

func getSpecRole(testSpec *api.CheTestsSpec) *v1.Role {
	return &v1.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Role",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      testHarnessRoleName,
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

func getRoleBindingSpec(testSpec *api.CheTestsSpec) *v1.RoleBinding {
	return &v1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testHarnessRoleBindingName,
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
			Name:     testHarnessRoleName,
		},
	}
}
