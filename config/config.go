package config

import (
	errHandler "k-bench/errHandler"
	"os"
	"path/filepath"

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
	Kubeconfig string
}

type deployStruct struct {
	Manifests []ManifestStruct
}

type ManifestStruct struct {
	Path      string
	Namespace string
	Cleanup   bool
}

// Returns a new k8s config instance built from the held .kubeconfig file.
func (c *connStruct) NewK8sConfig() (k8sConfig *rest.Config, err error) {
	k8sConfig, err = clientcmd.BuildConfigFromFlags("", c.Kubeconfig)
	if err != nil {
		return nil, errHandler.Error("error creating k8s config from kubeconfig", err)
	}

	return k8sConfig, nil
}

// Returns a new k8s client from the .kubeconfig file pointed to by the config struct.
func (c *connStruct) NewK8sClientSet() (k8sClientSet *kubernetes.Clientset, err error) {
	k8sConfig, err := c.NewK8sConfig()
	if err != nil {
		return nil, err
	}

	k8sClientSet, err = kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, errHandler.Error("error creating k8s client from kubeconfig", err)
	}

	if err != nil {
		return nil, err
	}
	return k8sClientSet, nil
}

// Singleton is used to prevent repeated file reads.
var configSingleton ConfigStruct

// Reads the passed *.yaml file.
// Returns a ConfigFile struct containing the contents.
func ReadConfig(path string) (*ConfigStruct, error) {
	var err error
	ext := filepath.Ext(path)
	log.Println(ext)
	if ext != ".yaml" && ext != ".yml" {
		return nil, errHandler.Error("invlid file extension, config fle must be a yaml file")
	}

	var bytes []byte
	bytes, err = os.ReadFile(path)
	if err != nil {
		return nil, errHandler.Error("error reading config file", err)
	}

	err = yaml.Unmarshal(bytes, &configSingleton)
	if err != nil {
		return nil, errHandler.Error("error parsing config file", err)
	}

	log.Debugf("read and parsed config file %s", path)
	return &configSingleton, nil
}

// Returns the underlying config singleton.
// ReadConfig() must have been called before Get()
func Get() *ConfigStruct {
	return &configSingleton
}
