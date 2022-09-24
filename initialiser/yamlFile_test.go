package initialiser

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestGetYamlFilepath(t *testing.T) {
	initer := New()
	expectedFilepath := "wibble.yaml"
	initer.yamlFilepath = expectedFilepath

	gotFilepath := initer.GetYamlFilepath()
	if gotFilepath != expectedFilepath {
		t.Fatalf("GetYamlFilepath returned unexpected value. Got %s. Expected %s.", gotFilepath, expectedFilepath)
	}
}

func TestValidateYamlFileValid(t *testing.T) {
	var validFiles []*os.File
	defer func() {
		err := cleanupFiles(validFiles)
		if err != nil {
			t.Fatal(err)
		}
	}()

	aFile, err := os.CreateTemp(os.TempDir(), "test-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	validFiles = append(validFiles, aFile)

	aFile, err = os.CreateTemp(os.TempDir(), "test-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	validFiles = append(validFiles, aFile)

	for i := 0; i < len(validFiles); i++ {
		initer := New()
		initer.yamlFilepath = validFiles[i].Name()
		err = initer.validateYamlFile()

		if err != nil {
			t.Fatalf("validateYamlFile() returned unexpected error for valid filepath %s: %v.", initer.yamlFilepath, err)
		}
	}

}

func TestValidateYamlFileValidLocal(t *testing.T) {
	localTmpDir, err := os.MkdirTemp(".", "tmp-*")
	log.Print(localTmpDir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.Remove(localTmpDir)
		if err != nil {
			t.Fatal(err)
		}
	}()

	var validFiles []*os.File
	defer func() {
		err := cleanupFiles(validFiles)
		if err != nil {
			t.Fatal(err)
		}
	}()

	aFile, err := os.CreateTemp(localTmpDir, "test-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	validFiles = append(validFiles, aFile)

	for i := 0; i < len(validFiles); i++ {
		initer := New()
		initer.yamlFilepath = validFiles[i].Name()
		err = initer.validateYamlFile()

		if err != nil {
			t.Fatalf("validateYamlFile() returned unexpected error for valid filepath %s: %v.", initer.yamlFilepath, err)
		}
	}
}

func TestValidateYamlFileBlank(t *testing.T) {
	initer := New()
	initer.yamlFilepath = ""
	err := initer.validateYamlFile()

	expectedErr := "blank filepath value"
	if err.Error() != expectedErr {
		t.Fatalf("validateYamlFile() returned unexpected error for blank yaml filepath. Got %s. Expected %s", err.Error(), expectedErr)
	}
}

func TestValidateYamlFileFlagValue(t *testing.T) {
	initer := New()
	initer.yamlFilepath = "-t"
	err := initer.validateYamlFile()

	expectedErr := "missing filepath value"
	if err.Error() != expectedErr {
		t.Fatalf("validateYamlFile() returned unexpected error for flag yaml filepath value. Got %s. Expected %s", err.Error(), expectedErr)
	}
}

func TestValidateYamlFileWrongExtension(t *testing.T) {
	initer := New()
	initer.yamlFilepath = "wibble.js"
	err := initer.validateYamlFile()

	expectedErr := "invalid filepath 'wibble.js'. Must be either *.yaml or *.yml"
	if err.Error() != expectedErr {
		t.Fatalf("validateYamlFile() returned unexpected error for non-yaml filepath value. Got %s. Expected %s", err.Error(), expectedErr)
	}
}

func TestValidateYamlFileNoFile(t *testing.T) {
	initer := New()
	initer.yamlFilepath = "missing.yaml"
	err := initer.validateYamlFile()

	expectedErr := fmt.Sprintf("stat %s: no such file or directory", initer.yamlFilepath)
	if err.Error() != expectedErr {
		t.Fatalf("validateYamlFile() returned unexpected error for missing file. Got %s. Expected %s", err.Error(), expectedErr)
	}
}

func cleanupFiles(files []*os.File) error {
	for i := 0; i < len(files); i++ {
		filepath := files[i].Name()
		files[i].Close()
		err := os.Remove(filepath)
		if err != nil {
			return err
		}
	}
	return nil
}
