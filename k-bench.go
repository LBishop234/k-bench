package main

import (
	"k-bench/configHandler"
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

	config, err := configHandler.Get(initer.GetYamlFilepath())
	if err != nil {
		log.WithField("error", err).Fatal("failed to read config file")
	}

	log.Print(config)
}
