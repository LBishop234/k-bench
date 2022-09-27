package configFile

import (
	errHandler "k-bench/errHandler"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type connStruct struct {
	Kubeconfig string
}

// Creates a k8s client set from the kubeconfig filepath.
func (c *connStruct) GetClientSet() (clientSet *kubernetes.Clientset, err error) {
	config, err := clientcmd.BuildConfigFromFlags("", c.Kubeconfig)
	if err != nil {
		return nil, errHandler.Error("error creating k8s config from kubeconfig", err)
	}

	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errHandler.Error("error creating k8s client from kubeconfig", err)
	}

	return clientSet, nil
}
