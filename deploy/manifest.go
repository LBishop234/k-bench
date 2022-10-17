package deploy

import (
	"context"
	"fmt"
	"k-bench/config"
	errhandler "k-bench/errHandler"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// Takes a path to a manifest file or directory containing manifest files.
// Applies all found files to the cluster.
func DeployManifests(path, namespace string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return errhandler.Error("failed to convert file relative path to full path", err)
	}
	filePaths, err := findManifests(path)
	if err != nil {
		return errhandler.Error("failed to handle manifest file(s)", err)
	}

	for i := 0; i < len(filePaths); i++ {
		err = applyManifest(filePaths[i], namespace)
		if err != nil {
			return errhandler.Error("failed to apply manifest to cluster")
		}
	}
	return nil
}

// Takes a path to a manifest file or directory containing manifest files.
// Removes all found files from the cluster.
func RemoveManfiests(path, namespace string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return errhandler.Error("failed to convert file relative path to full path", err)
	}
	filePaths, err := findManifests(path)
	if err != nil {
		return errhandler.Error("failed to handle manifest file(s)", err)
	}

	for i := 0; i < len(filePaths); i++ {
		err = removeManifest(filePaths[i], namespace)
		if err != nil {
			return errhandler.Error("failed to remove manifest from cluster")
		}
	}
	return nil
}

// Checks if a path points to a file or directory.
// If a directory is found, discovers all files in the directory.
// Validates all filepaths.
// Returns a list of filepaths.
func findManifests(path string) (filepaths []string, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return filepaths, errhandler.Error("error checking manifest file/directory", err)
	}

	if fileInfo.IsDir() {
		dirFiles, err := os.ReadDir(path)
		if err != nil {
			return filepaths, errhandler.Error("error checking manifest directory contents", err)
		}
		for i := 0; i < len(dirFiles); i++ {
			fullPath := fmt.Sprintf("%s/%s", path, dirFiles[i].Name())
			subPaths, err := findManifests(fullPath)
			if err != nil {
				return filepaths, errhandler.Error("error handling manifest directory contents", err)
			}
			filepaths = append(filepaths, subPaths...)
		}
	} else {
		filepaths = append(filepaths, path)
	}

	return filepaths, nil
}

// Applys a manifest file to the cluster.
func applyManifest(filepath string, namespace string) error {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return errhandler.Error("error opening manifest file", err)
	}

	var manifest unstructured.Unstructured
	err = yaml.Unmarshal(file, &manifest.Object)
	if err != nil {
		return errhandler.Error("error unmarshalling manifest yaml", err)
	}

	aConfig := config.Get()
	k8sConfig, err := aConfig.Connection.NewK8sConfig()
	if err != nil {
		return errhandler.Error("", err)
	}
	k8sClient, err := dynamic.NewForConfig(k8sConfig)
	if err != nil {
		return errhandler.Error("error creating dynamic client connection", err)
	}

	var manifestType schema.GroupVersionResource
	switch manifest.GetKind() {
	case "Deployment":
		manifestType = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	case "Service":
		manifestType = schema.GroupVersionResource{Group: "core", Version: "v1", Resource: "services"}
	}

	result, err := k8sClient.Resource(manifestType).Namespace(namespace).Create(context.TODO(), &manifest, metav1.CreateOptions{})
	if err != nil {
		return errhandler.Error("error applying manifest", err)
	}

	log.Debugf("Applied %s to cluster", result.GetName())
	return nil
}

// Removes a manifest file from the cluster.
func removeManifest(filepath string, namespace string) error {
	return nil
}
