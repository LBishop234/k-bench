package configHandler

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type configFile struct {
	Connection connStruct
}

// Reads the passed *.yaml file.
// Returns a ConfigFile struct containing the contents.
func readConfigFile(filepath string) (aFile configFile, err error) {
	logFields := log.Fields{
		"filepath": filepath,
	}

	bytes, err := os.ReadFile(filepath)
	if err != nil {
		log.WithFields(logFields).WithField("error", err).Debug("error reading config file %s", filepath)
		return aFile, err
	}

	err = yaml.Unmarshal(bytes, &aFile)
	if err != nil {
		log.WithFields(logFields).WithField("error", err).Debug("error unmarshalling config file %s", filepath)
		return aFile, err
	}

	log.WithFields(logFields).Debugf("read and parsed config file %s", filepath)
	return aFile, nil
}
