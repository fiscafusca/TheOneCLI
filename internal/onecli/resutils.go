package onecli

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8syaml "sigs.k8s.io/yaml"
)

type Resource struct {
	GroupVersionKind *schema.GroupVersionKind
	Object           unstructured.Unstructured
}

// NewResourcesFromFiles creates new deployable resources from the YAML manifests
func NewResourcesFromFiles(resourcesPath string) ([]Resource, error) {
	var stream []byte
	var err error
	var resources []Resource

	files, err := os.ReadDir(resourcesPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filePath := path.Join(resourcesPath, file.Name())
		stream, err = os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		r, err := createResourcesFromBuffer(stream, resourcesPath)
		if err != nil {
			return nil, err
		}
		resources = append(resources, r...)
	}

	return resources, nil
}

// createResourcesFromBuffer creates new deployable resources from a byte stream
func createResourcesFromBuffer(stream []byte, filepath string) ([]Resource, error) {
	var resources []Resource

	re := regexp.MustCompile(`\n---\n`)
	for _, resourceYAML := range re.Split(string(stream), -1) {
		if len(resourceYAML) == 0 {
			continue
		}

		u := unstructured.Unstructured{Object: map[string]interface{}{}}
		if err := k8syaml.Unmarshal([]byte(resourceYAML), &u.Object); err != nil {
			return nil, fmt.Errorf("resource %s: %s", filepath, err)
		}
		gvk := u.GroupVersionKind()

		resource := Resource{
			GroupVersionKind: &gvk,
			Object:           u,
		}

		resources = append(resources, resource)
	}

	return resources, nil
}
