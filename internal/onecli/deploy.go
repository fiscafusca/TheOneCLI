package onecli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

var gvrNamespaces = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}

func NewDeployCommand() *cobra.Command {
	initCmd := &cobra.Command{

		Use:   "deploy",
		Short: "Deploy the resources",
		Long: `A veeeeeeeeeeeeeeeeeeeeery
loooooooooooooooooooooooooooooooooooong
descriptiooooooooooooooooooooooooooooon.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Println("Deploying...")

			cfg, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
			if err != nil {
				return err
			}

			clients := &K8sClients{
				dynamic:   dynamic.NewForConfigOrDie(cfg),
				discovery: discovery.NewDiscoveryClientForConfigOrDie(cfg),
			}

			resources, err := NewResourcesFromFiles("./resources")
			if err != nil {
				return err
			}
			return Deploy(clients, opts.namespace, resources)
		},
	}

	initCmd.Flags().StringVar(&opts.namespace, "namespace", "default", "namespace")
	return initCmd
}

// Deploy ensures the existance of the namespace and calls the apply function for each resource
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

// apply is the function that actually creates/patches resources on the cluster
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
