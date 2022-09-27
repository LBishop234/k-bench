package configHandler

import (
	errHandler "k-bench/errHandler"
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
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return aFile, errHandler.Error("error reading config file", err)
	}

	err = yaml.Unmarshal(bytes, &aFile)
	if err != nil {
		return aFile, errHandler.Error("error unmarshalling config file", err)
	}

	log.Debugf("read and parsed config file %s", filepath)
	return aFile, nil
}
