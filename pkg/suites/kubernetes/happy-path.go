package kubernetes

import (
	"context"
	"fmt"
	util2 "github.com/flacatus/che-inspector/pkg/util"

	"github.com/flacatus/che-inspector/pkg/api"
	"github.com/flacatus/che-inspector/pkg/common/client"
	"github.com/flacatus/che-inspector/pkg/common/clog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Comment
func DeployHappyPath(k8sClient *client.K8sClient, testSpec *api.CheTestsSpec) (err error) {
	if _, err := k8sClient.Kube().CoreV1().Namespaces().Get(context.TODO(),testSpec.Namespace, metav1.GetOptions{}); err != nil {
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

	terminated, err := waitForContainerToBeTerminated(k8sClient, testSpec, pod.Name)
	if terminated {
		err = util2.CopyArtifactsFromPod(testSpec.Artifacts.FromContainerPath, testSpec.Artifacts.To, pod.Name, testSpec.Namespace, artifactsDownloadContainerName)
	} else {
		return fmt.Errorf("Failed to get test pod status")
	}

	return err
}

// GetNamespaceSpec return namespace object
func GetNamespaceSpec(testSpec *api.CheTestsSpec) *v1.Namespace {
	return &v1.Namespace{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: testSpec.Namespace,
		},
	}
}
