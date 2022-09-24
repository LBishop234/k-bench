package initialiser

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/sirupsen/logrus"
)

type Initialiser struct {
	yamlFilepath string
	traceMode    bool
	quietMode    bool
}

func (i *Initialiser) Log() {
	fmt.Println(i)
}

// Returns a new Initialiser struct.
func New() Initialiser {
	return Initialiser{}
}

// Sets the Initialiser values from command line arguments.
func (i *Initialiser) SetFromArgs() {
	flag.BoolVar(&i.traceMode, "t", false, "Trace mode. Enables trace logs. Not recommended for non debugging use")
	flag.BoolVar(&i.quietMode, "q", false, "Quiet mode. Reduces the number of output logs")

	flag.Parse()

	i.yamlFilepath = os.Args[len(os.Args)-1]
}

// Returns the yaml filepath this.
func (i *Initialiser) GetYamlFilepath() string {
	return i.yamlFilepath
}

// Initialise any executable dependencies.
func (i *Initialiser) Initialise() error {
	i.setLogrusFormat()

	err := i.setLogrusLevel()
	if err != nil {
		return err
	}

	err = i.validateYamlFile()
	if err != nil {
		return err
	}

	return nil
}

// Sets the Logrus Formatter
func (i *Initialiser) setLogrusFormat() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
}

// Sets the underlying Logrus level according to Quiet Mode and Trace Mode flags.
func (i *Initialiser) setLogrusLevel() error {
	if i.quietMode && i.traceMode {
		return fmt.Errorf("both Quiet Mode and Trace Mode cannot be set simultaneously")
	}

	if i.quietMode {
		logrus.SetLevel(logrus.InfoLevel)
	}

	if i.traceMode {
		logrus.SetLevel(logrus.TraceLevel)
	}

	if !i.quietMode && !i.traceMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	return nil
}

// Validates the passed yaml filepath value.
func (i *Initialiser) validateYamlFile() error {
	if i.yamlFilepath == "" {
		return fmt.Errorf("blank filepath value")
	}

	match, err := regexp.MatchString("-([a-z]|[A-Z])", i.yamlFilepath)
	if err != nil {
		return err
	}
	if match {
		return fmt.Errorf("missing filepath value")
	}

	fType := filepath.Ext(i.yamlFilepath)
	if fType == "" || (fType != ".yaml" && fType != ".yml") {
		return fmt.Errorf("invalid filepath. Must be either *.yaml or *.yml")
	}

	_, err = os.Stat(i.yamlFilepath)
	if err != nil {
		return err
	}

	return nil
}
