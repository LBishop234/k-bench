package initialiser

import (
	"flag"
	"fmt"
	errHandler "k-bench/errHandler"
	"os"
	"path/filepath"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type Initialiser struct {
	yamlFilepath string
	debugMode    bool
	quietMode    bool
}

// Returns a new Initialiser struct.
func New() Initialiser {
	return Initialiser{}
}

// Sets the Initialiser values from command line arguments.
func (i *Initialiser) SetFromArgs() {
	flag.BoolVar(&i.debugMode, "d", false, "Debug mode. Enables debug logs. Not recommended for non debugging use.")
	flag.BoolVar(&i.quietMode, "q", false, "Quiet mode. Reduces the number of output logs.")
	// Flag Idea: output -o. Outputs key info to file, i.e. k6 outputs and other metrics.

	flag.Parse()

	i.yamlFilepath = os.Args[len(os.Args)-1]
}

// Returns the yaml filepath this.
func (i *Initialiser) GetYamlFilepath() string {
	return i.yamlFilepath
}

// Initialise any executable dependencies.
func (i *Initialiser) Initialise() error {
	err := i.initLogrus()
	if err != nil {
		return errHandler.Error("error initalising logrus", err)
	}

	err = i.validateYamlFile()
	if err != nil {
		return errHandler.Error("error validating yaml config file", err)
	}

	return nil
}

// Initialises logrus.
// Sets the underlying Logrus level according to Quiet Mode and Trace Mode flags.
// If neither flag is set, defaults to InfoLevel.
func (i *Initialiser) initLogrus() error {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		PadLevelText:     true,
	})

	if i.quietMode && i.debugMode {
		errHandler.Error("both Quiet Mode and Debug Mode cannot be set simultaneously")
	}

	if i.quietMode {
		log.SetLevel(log.WarnLevel)
	}

	if i.debugMode {
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}

	if !i.quietMode && !i.debugMode {
		log.SetLevel(log.InfoLevel)
	}

	log.Debugf("logrus initalised with log level %s", log.GetLevel())
	return nil
}

// Validates the passed yaml filepath value.
func (i *Initialiser) validateYamlFile() error {
	if i.yamlFilepath == "" {
		return errHandler.Error("initialiser got invalid yaml file", fmt.Errorf("blank yaml filepath value"))
	}

	match, err := regexp.MatchString("-([a-z]|[A-Z])", i.yamlFilepath)
	if err != nil {
		return errHandler.Error("error checking yaml file regex", err)
	}
	if match {
		return errHandler.Error(fmt.Sprintf("invalid yaml filepath value %s", i.yamlFilepath))
	}

	fType := filepath.Ext(i.yamlFilepath)
	if fType == "" || (fType != ".yaml" && fType != ".yml") {
		return errHandler.Error("initialiser got invalid yaml file", fmt.Errorf("invalid filepath '%s'. Must be either *.yaml or *.yml", i.yamlFilepath))
	}

	_, err = os.Stat(i.yamlFilepath)
	if err != nil {
		return errHandler.Error("yaml file does not exist", err)
	}

	log.WithField("filepath", i.yamlFilepath).Debug("initialiser got valid yaml file")
	return nil
}
