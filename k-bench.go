package main

import (
	"k-bench/initialiser"

	"github.com/sirupsen/logrus"
)

func main() {
	initer := initialiser.New()
	initer.SetFromArgs()
	err := initer.Initialise()
	if err != nil {
		logrus.WithField("error", err).Fatal("failed to initialise correctly")
	}

	initer.Log()
}
