package configHandler

import (
	errHandler "k-bench/errHandler"

	log "github.com/sirupsen/logrus"
)

type Config struct{}

// Reads the config file passed.
// Returns the config struct
func Get(filepath string) (conf Config, err error) {
	aFile, err := readConfigFile(filepath)
	if err != nil {
		return conf, errHandler.Error("error reading config file", err)
	}
	log.Print(aFile)

	return Config{}, nil
}
