package configHandler

import (
	"k-bench/configHandler/configFile"
	errHandler "k-bench/errHandler"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	k8sClient *kubernetes.Clientset
}

// Reads the config file passed.
// Returns the config struct
func Get(filepath string) (conf Config, err error) {
	aFile, err := configFile.ReadConfigFile(filepath)
	if err != nil {
		return conf, errHandler.Error("error reading config file", err)
	}
	log.Print(aFile)

	conf.k8sClient, err = aFile.Connection.GetClientSet()
	if err != nil {
		return conf, errHandler.Error("error creating kubernetes client", err)
	}

	return conf, nil
}
