package onecli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var gvrNamespaces = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}

func NewDeployCommand() *cobra.Command {
	initCmd := &cobra.Command{

		Use:   "deploy [CONTEXT]",
		Short: "Deploy the resources",
		Long: `A veeeeeeeeeeeeeeeeeeeeery
loooooooooooooooooooooooooooooooooooong
descriptiooooooooooooooooooooooooooooon.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Deploying...")

			var err error
			var clients *K8sClients

			contextsMap := viper.Get("contexts")

			// deploy all contexts
			if len(args) == 0 {
				for context := range contextsMap.(map[string]interface{}) {
					clients, err = createK8sClients(context)
					if err != nil {
						return err
					}
					err = deployInContext(context, clients, contextsMap.(map[string]interface{}))
					if err != nil {
						return err
					}
				}
			} else {
				// deploy in given context
				clients, err = createK8sClients(args[0])
				if err != nil {
					return err
				}
				err = deployInContext(args[0], clients, contextsMap.(map[string]interface{}))
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	initCmd.Flags().StringVar(&opts.namespace, "namespace", "", "namespace")
	return initCmd
}

// deployInContext deploys only the resources in the config file in the selected context
func deployInContext(context string, clients *K8sClients, contextsMap map[string]interface{}) error {
	for _, res := range contextsMap[context].([]interface{}) {
		resources, err := NewResourcesFromFiles(fmt.Sprint(res))
		if err != nil {
			return err
		}
		err = deployResources(clients, opts.namespace, resources)
		if err != nil {
			return err
		}
	}
	return nil
}

// deployResources ensures the existance of the namespace and calls the apply function for each resource
func deployResources(clients *K8sClients, namespace string, resources []Resource) error {
	if namespace != "" {
		if err := ensureNamespaceExistence(clients, namespace); err != nil {
			return fmt.Errorf("error ensuring namespace existence for namespace %s: %w", namespace, err)
		}
	}

	// apply the resources
	for _, res := range resources {
		err := applyResource(clients, res)
		if err != nil {
			return fmt.Errorf("error applying resource %+v: %w", res, err)
		}
	}
	return nil
}

// applyResource is the function that actually creates/patches resources on the cluster
func applyResource(clients *K8sClients, res Resource) error {
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
