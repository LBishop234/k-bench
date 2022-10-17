package main

import (
	"k-bench/config"
	"k-bench/deploy"
	errhandler "k-bench/errHandler"
	"k-bench/initialiser"

	log "github.com/sirupsen/logrus"
)

func main() {
	initer := initialiser.New()
	initer.SetFromArgs()
	err := initer.Initialise()
	if err != nil {
		log.WithField("error", err).Fatal("failed to initialise correctly")
	}
	log.Info("initialised k-bench")

	aConfig, err := config.ReadConfig(initer.GetYamlFilepath())
	if err != nil {
		log.WithField("error", err).Fatal("failed to read and parse config file")
	}
	log.Info("read and parsed config file")

	err = deployManifests(aConfig.Deploy.Manifests)
	if err != nil {
		log.WithField("error", err).Fatal("failed to deploy all manifests")
	}
	log.Info("deployed all manifests to cluster")

	defer func() {
		err = cleanupManifests(aConfig.Deploy.Manifests)
		if err != nil {
			log.WithField("error", err).Fatal("failed to cleanup marked manifests")
		}
		log.Info("Cleaned up marked manifests")
	}()
}

// deploys passed manifests to the cluster
func deployManifests(manifests []config.ManifestStruct) error {
	var deployErrs []error
	for i := 0; i < len(manifests); i++ {
		err := deploy.DeployManifests(manifests[i].Path, manifests[i].Namespace)
		if err != nil {
			deployErrs = append(deployErrs, err)
		}
	}
	if len(deployErrs) > 0 {
		return errhandler.Error("error deploying manifests", deployErrs...)
	}
	return nil
}

// removes any manifests flagged for cleanup from the cluster.
func cleanupManifests(manifests []config.ManifestStruct) error {
	var cleanupErrs []error
	for i := 0; i < len(manifests); i++ {
		if manifests[i].Cleanup {
			err := deploy.RemoveManfiests(manifests[i].Path, manifests[i].Namespace)
			if err != nil {
				cleanupErrs = append(cleanupErrs, err)
			}
		}
	}
	if len(cleanupErrs) > 0 {
		return errhandler.Error("error removing manifests", cleanupErrs...)
	}
	return nil
}
