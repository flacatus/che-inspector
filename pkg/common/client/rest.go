package client

import (
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Comment
func InitRestClient() (*rest.Config, error, *corev1client.CoreV1Client) {
	// Instantiate loader for kubeconfig file.
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Get a rest.Config from the kubeConfig file.  This will be passed into all
	// the client objects we create.
	restConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	// Create a Kubernetes core/v1 client.
	coreClient, err := corev1client.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	return restConfig, err, coreClient
}
