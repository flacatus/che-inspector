package kubernetes

import (
	"fmt"
	"github.com/flacatus/che-inspector/pkg/common/client"
	"github.com/flacatus/che-inspector/pkg/common/instance"
	"github.com/flacatus/che-inspector/pkg/common/util"
)

func DeployHappyPath(k8sClient *client.K8sClient, testSpec *instance.CheTestsSpec) (err error) {
	pod, err := k8sClient.Kube().CoreV1().Pods(testSpec.Namespace).Create(GetTestSuitePodSpec(testSpec))
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
