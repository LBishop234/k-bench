package main

import (
	"k-bench/config"
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

	config, err := config.ReadConfigFile(initer.GetYamlFilepath())
	if err != nil {
		log.WithField("error", err).Fatal("failed to read and parse config file")
	}
	log.Info("parsed config file")

	log.Print(config)
}
