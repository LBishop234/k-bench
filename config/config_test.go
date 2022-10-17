package config

import (
	"fmt"
	errHandler "k-bench/errHandler"
	"testing"
)

func TestReadConfig(t *testing.T) {
	path := "./test-resources/test.yaml"

	aFile, err := ReadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	xpctKubeConfig := "/Users/admin/.kube/config"
	if aFile.Connection.Kubeconfig != xpctKubeConfig {
		t.Fatalf("unexpected connection.kubeconfig value. Got '%s'. Expected '%s'.", aFile.Connection.Kubeconfig, xpctKubeConfig)
	}

	xpctdManifestPath := "/test/wibble.yaml"
	xpctdManifestNamespace := "default"
	xpctdManifestCleanup := true
	if aFile.Deploy.Manifests[0].Path != xpctdManifestPath {
		t.Fatalf("unexpected deploy.manifests.path value. Got '%s'. Expected '%s'.", aFile.Deploy.Manifests[0].Path, xpctdManifestPath)
	}
	if aFile.Deploy.Manifests[0].Namespace != xpctdManifestNamespace {
		t.Fatalf("unexpected deploy.manifests.namespace value. Got '%s'. Expected '%s'.", aFile.Deploy.Manifests[0].Namespace, xpctdManifestNamespace)
	}
	if aFile.Deploy.Manifests[0].Cleanup != xpctdManifestCleanup {
		t.Fatalf("unexpected deploy.manifests.cleanup value. Got '%t'. Expected '%t'.", aFile.Deploy.Manifests[0].Cleanup, xpctdManifestCleanup)
	}

	xpctdManifestPath = "/test/dir/"
	xpctdManifestNamespace = "testing"
	xpctdManifestCleanup = false
	if aFile.Deploy.Manifests[1].Path != xpctdManifestPath {
		t.Fatalf("unexpected deploy.manifests.path value. Got '%s'. Expected '%s'.", aFile.Deploy.Manifests[1].Path, xpctdManifestPath)
	}
	if aFile.Deploy.Manifests[1].Namespace != xpctdManifestNamespace {
		t.Fatalf("unexpected deploy.manifests.namespace value. Got '%s'. Expected '%s'.", aFile.Deploy.Manifests[1].Namespace, xpctdManifestNamespace)
	}
	if aFile.Deploy.Manifests[1].Cleanup != xpctdManifestCleanup {
		t.Fatalf("unexpected deploy.manifests.cleanup value. Got '%t'. Expected '%t'.", aFile.Deploy.Manifests[1].Cleanup, xpctdManifestCleanup)
	}
}

func TestReadConfigFileBadFilepath(t *testing.T) {
	_, err := ReadConfig("./test-resources/missing.yaml")
	xpctdErr := errHandler.Error("error reading config file", fmt.Errorf("open ./test-resources/missing.yaml: no such file or directory"))
	if err.Error() != xpctdErr.Error() {
		t.Fatalf("unexpected error from ReadConfig() with bad filepath. Got '%v'. Expected '%v'.", err, xpctdErr)
	}
}

func TestReadConfigFileBadFileType(t *testing.T) {
	_, err := ReadConfig("./test-resources/test.json")
	xpctdErr := errHandler.Error("invlid file extension, config fle must be a yaml file")
	if err.Error() != xpctdErr.Error() {
		t.Fatalf("unexpected error from ReadConfig() with bad file type. Got '%v'. Expected '%v'.", err, xpctdErr)
	}
}
