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
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Deploying...")

			var err error
			var clients *K8sClients

			contextsMap := viper.Get("contexts")

			if contextsMap == nil {
				// deploy in current context
				clients, err = createK8sClients("")
				if err != nil {
					return err
				}
				resources, err := NewResourcesFromFiles("./resources")
				if err != nil {
					return err
				}
				return deploy(clients, opts.namespace, resources)
			}
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

	initCmd.Flags().StringVar(&opts.namespace, "namespace", "default", "namespace")
	return initCmd
}

func deployInContext(context string, clients *K8sClients, contextsMap map[string]interface{}) error {
	for _, res := range contextsMap[context].([]interface{}) {
		resources, err := NewResourcesFromFiles(fmt.Sprint(res))
		if err != nil {
			return err
		}
		err = deploy(clients, opts.namespace, resources)
		if err != nil {
			return err
		}
	}
	return nil
}

// deploy ensures the existance of the namespace and calls the apply function for each resource
func deploy(clients *K8sClients, namespace string, resources []Resource) error {
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
