package initialiser

import (
	"flag"
	"fmt"
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
		return err
	}

	err = i.validateYamlFile()
	if err != nil {
		return err
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
		return fmt.Errorf("both Quiet Mode and Debug Mode cannot be set simultaneously")
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
	logFields := log.Fields{
		"filepath": i.yamlFilepath,
	}

	if i.yamlFilepath == "" {
		err := fmt.Errorf("blank filepath value")
		log.WithFields(logFields).WithField("error", err).Debug("initialiser got invalid yaml file")
		return err
	}

	match, err := regexp.MatchString("-([a-z]|[A-Z])", i.yamlFilepath)
	if err != nil {
		log.WithFields(logFields).WithField("error", err).Debug("error checking yaml file regex")
		return err
	}
	if match {
		err = fmt.Errorf("missing filepath value")
		log.WithFields(logFields).WithField("error", err).Debug("initialiser got invalid yaml file")
		return err
	}

	fType := filepath.Ext(i.yamlFilepath)
	if fType == "" || (fType != ".yaml" && fType != ".yml") {
		err = fmt.Errorf("invalid filepath '%s'. Must be either *.yaml or *.yml", i.yamlFilepath)
		log.WithFields(logFields).WithField("error", err).Debug("initialiser got invalid yaml file")
		return err
	}

	_, err = os.Stat(i.yamlFilepath)
	if err != nil {
		log.WithFields(logFields).WithField("error", err).Debug("initialiser got invalid yaml file")
		return err
	}

	log.WithFields(logFields).Debug("initialiser got valid yaml file")
	return nil
}
