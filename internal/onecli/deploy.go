package onecli

import (
	"context"
	"fmt"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/jsonmergepatch"
	"k8s.io/apimachinery/pkg/util/mergepatch"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
)

var gvrNamespaces = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}

func NewDeployCommand() *cobra.Command {
	initCmd := &cobra.Command{

		Use:   "deploy",
		Short: "Deploy the resources",
		Long: `A veeeeeeeeeeeeeeeeeeeeery
loooooooooooooooooooooooooooooooooooong
descriptiooooooooooooooooooooooooooooon.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Deploying...")

			cfg, err := clientcmd.BuildConfigFromFlags("", "/home/giorgia/.kube/config")
			if err != nil {
				return err
			}

			clients := &K8sClients{
				dynamic:   dynamic.NewForConfigOrDie(cfg),
				discovery: discovery.NewDiscoveryClientForConfigOrDie(cfg),
			}
			path := path.Join(args[0], "pod.yaml")
			resources, err := NewResourcesFromFile(path)
			if err != nil {
				return err
			}
			return Deploy(clients, flags.namespace, resources)
		},
	}

	initCmd.Flags().StringVar(&flags.namespace, "namespace", "default", "namespace")
	return initCmd
}

func Deploy(clients *K8sClients, namespace string, resources []Resource) error {
	if namespace != "" {
		if err := ensureNamespaceExistence(clients, namespace); err != nil {
			return fmt.Errorf("error ensuring namespace existence for namespace %s: %w", namespace, err)
		}
	}

	// apply the resources
	for _, res := range resources {
		err := apply(clients, res)
		if err != nil {
			return fmt.Errorf("error applying resource %+v: %w", res, err)
		}
	}
	return nil
}

func apply(clients *K8sClients, res Resource) error {
	gvr, err := FromGVKtoGVR(clients.discovery, res.Object.GroupVersionKind())
	if err != nil {
		return err
	}

	onClusterObj, err := GetResource(gvr, clients, res)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return CreateResource(gvr, clients, res)
		}

		return err
	}

	return PatchResource(gvr, clients, res, onClusterObj)
}

func CreateResource(gvr schema.GroupVersionResource, clients *K8sClients, res Resource) error {
	fmt.Printf("Creating %s: %s\n", res.Object.GetKind(), res.Object.GetName())

	// creates kubectl.kubernetes.io/last-applied-configuration annotation
	// inside the resource except for Secrets, ConfigMaps, and CRDs
	originAnn := res.Object.GetAnnotations()
	if originAnn == nil {
		originAnn = make(map[string]string)
	}
	objJSON, err := res.Object.MarshalJSON()
	if err != nil {
		return err
	}
	originAnn[corev1.LastAppliedConfigAnnotation] = string(objJSON)
	res.Object.SetAnnotations(originAnn)

	// var resourceInterface dynamic.ResourceInterface
	// switch res.Namespaced {
	// case True:
	// 	resourceInterface = clients.dynamic.Resource(gvr).Namespace(res.Object.GetNamespace())
	// case False:
	// 	resourceInterface = clients.dynamic.Resource(gvr)
	// case None:
	// 	return fmt.Errorf("resource %s %s is unknown in cluster, can't apply it", res.Object.GetName(), gvr)
	// }

	resourceInterface := clients.dynamic.Resource(gvr).Namespace(res.Object.GetNamespace())

	_, err = resourceInterface.Create(context.Background(), &res.Object, metav1.CreateOptions{})
	return err
}

func GetResource(gvr schema.GroupVersionResource, clients *K8sClients, res Resource) (*unstructured.Unstructured, error) {
	return clients.dynamic.Resource(gvr).
		Namespace(res.Object.GetNamespace()).
		Get(context.Background(), res.Object.GetName(), metav1.GetOptions{})
}

func PatchResource(gvr schema.GroupVersionResource, clients *K8sClients, res Resource, onClusterObj *unstructured.Unstructured) error {
	// create the patch
	patch, patchType, err := createPatch(*onClusterObj, res)
	if err != nil {
		return errors.Wrap(err, "failed to create patch")
	}

	var resourceInterface dynamic.ResourceInterface
	switch res.Namespaced {
	case True:
		resourceInterface = clients.dynamic.Resource(gvr).Namespace(res.Object.GetNamespace())
	case False:
		resourceInterface = clients.dynamic.Resource(gvr)
	case None:
		return fmt.Errorf("resource %s %s is unknown in cluster, can't patch it", res.Object.GetName(), gvr)
	}

	_, err = resourceInterface.Patch(context.Background(), res.Object.GetName(), patchType, patch, metav1.PatchOptions{})
	return err
}

func createPatch(currentObj unstructured.Unstructured, target Resource) ([]byte, types.PatchType, error) {
	// Get the resource in the cluster
	currentJSON, err := currentObj.MarshalJSON()
	if err != nil {
		return nil, "", errors.Wrap(err, "serializing live configuration")
	}

	// Get last applied config from annotation if exists
	lastAppliedConfigAnnotation := ""
	lastAppliedConfigAnnotationFound := false
	var targetJSON []byte
	annotations := currentObj.GetAnnotations()
	if annotations != nil {
		lastAppliedConfigAnnotation, lastAppliedConfigAnnotationFound = annotations[corev1.LastAppliedConfigAnnotation]
	}

	if lastAppliedConfigAnnotationFound {
		annotatedTarget, err := annotateWithLastApplied(target)
		if err != nil {
			return nil, "", err
		}
		targetJSON, err = annotatedTarget.MarshalJSON()
		if err != nil {
			return nil, "", err
		}
	} else {
		targetJSON, err = target.Object.MarshalJSON()
		if err != nil {
			return nil, "", err
		}
	}

	versionedObject, err := scheme.Scheme.New(*target.GroupVersionKind)
	if err != nil && !runtime.IsNotRegisteredError(err) {
		return nil, "", err
	}

	// use a three way json merge if the resource is a CRD
	if runtime.IsNotRegisteredError(err) {
		// fall back to generic JSON merge patch
		patchType := types.MergePatchType
		preconditions := []mergepatch.PreconditionFunc{mergepatch.RequireKeyUnchanged("apiVersion"),
			mergepatch.RequireKeyUnchanged("kind"), mergepatch.RequireMetadataKeyUnchanged("name")}
		patch, err := jsonmergepatch.CreateThreeWayJSONMergePatch([]byte(lastAppliedConfigAnnotation), targetJSON, currentJSON, preconditions...)
		return patch, patchType, err
	}

	patchMeta, err := strategicpatch.NewPatchMetaFromStruct(versionedObject)
	if err != nil {
		return nil, types.StrategicMergePatchType, errors.Wrap(err, "unable to create patch metadata from object")
	}

	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(lastAppliedConfigAnnotation), targetJSON, currentJSON, patchMeta, true)
	return patch, types.StrategicMergePatchType, err
}

func ensureNamespaceExistence(clients *K8sClients, namespace string) error {
	ns := &unstructured.Unstructured{}
	ns.SetUnstructuredContent(map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Namespace",
		"metadata": map[string]interface{}{
			"name": namespace,
		},
	})

	fmt.Printf("Creating namespace %s\n", namespace)
	if _, err := clients.dynamic.Resource(gvrNamespaces).Create(context.Background(), ns, metav1.CreateOptions{}); err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func annotateWithLastApplied(res Resource) (unstructured.Unstructured, error) {
	annotatedRes := res.Object.DeepCopy()
	annotations := annotatedRes.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	if _, found := annotations[corev1.LastAppliedConfigAnnotation]; found {
		delete(annotations, corev1.LastAppliedConfigAnnotation)
		annotatedRes.SetAnnotations(annotations)
	}

	resJSON, err := annotatedRes.MarshalJSON()
	if err != nil {
		return unstructured.Unstructured{}, err
	}

	annotations[corev1.LastAppliedConfigAnnotation] = string(resJSON)
	annotatedRes.SetAnnotations(annotations)

	return *annotatedRes, nil
}
