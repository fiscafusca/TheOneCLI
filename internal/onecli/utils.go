package onecli

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	k8syaml "sigs.k8s.io/yaml"
)

type ResourceNamespace int8

const (
	None ResourceNamespace = iota
	True
	False
)

type Resource struct {
	Filepath         string
	GroupVersionKind *schema.GroupVersionKind
	Object           unstructured.Unstructured
	Namespaced       ResourceNamespace
}

type ResourceList struct {
	Kind      *schema.GroupVersionKind `json:"kind"`
	Resources []string                 `json:"resources"`
}

type K8sClients struct {
	dynamic   dynamic.Interface
	discovery discovery.DiscoveryInterface
}

// FromGVKtoGVR converts Group Version Kind to Group Version Resource
func FromGVKtoGVR(discoveryClient discovery.DiscoveryInterface, gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))
	a, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	return a.Resource, nil
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
			Filepath:         filepath,
			GroupVersionKind: &gvk,
			Object:           u,
			Namespaced:       None,
		}

		resources = append(resources, resource)
	}

	return resources, nil
}
