package config

import (
	errHandler "k-bench/errHandler"
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ConfigStruct struct {
	Connection connStruct
	Deploy     deployStruct
}

type connStruct struct {
	clientSetSync sync.Once
	clientSet     *kubernetes.Clientset
	Kubeconfig    string
}

type deployStruct struct {
	Manifests []manifestStruct
}

type manifestStruct struct {
	Path      string
	Namespace string
	Cleanup   bool
}

func (c *connStruct) NewK8sClientSet() (*kubernetes.Clientset, error) {
	var err error
	c.clientSetSync.Do(func() {
		var config *rest.Config
		config, err = clientcmd.BuildConfigFromFlags("", c.Kubeconfig)
		if err != nil {
			err = errHandler.Error("error creating k8s config from kubeconfig", err)
			return
		}

		c.clientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			err = errHandler.Error("error creating k8s client from kubeconfig", err)
			return
		}
	})

	if err != nil {
		return nil, err
	}
	return c.clientSet, nil
}

// Reads the passed *.yaml file.
// Returns a ConfigFile struct containing the contents.
func ReadConfigFile(path string) (aFile ConfigStruct, err error) {
	ext := filepath.Ext(path)
	log.Println(ext)
	if ext != ".yaml" && ext != ".yml" {
		return aFile, errHandler.Error("invlid file extension, config fle must be a yaml file")
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return aFile, errHandler.Error("error reading config file", err)
	}

	err = yaml.Unmarshal(bytes, &aFile)
	if err != nil {
		return aFile, errHandler.Error("error parsing config file", err)
	}

	log.Debugf("read and parsed config file %s", path)
	return aFile, nil
}
