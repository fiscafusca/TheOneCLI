package onecli

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Resource struct {
	GroupVersionKind *schema.GroupVersionKind
	Object           unstructured.Unstructured
}

// NewResourcesFromFiles creates new deployable resources from the YAML manifests
func NewResourcesFromFiles(resourcesPath string) ([]Resource, error) {

	// read streams from manifest files and call createResourcesFromBuffer

	return nil, nil
}

// createResourcesFromBuffer creates new deployable resources from a byte stream
func createResourcesFromBuffer(stream []byte, filepath string) ([]Resource, error) {

	// create k8s resources from a byte stream

	return nil, nil
}
