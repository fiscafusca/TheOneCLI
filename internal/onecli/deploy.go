package onecli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

			// create the k8s clients...
			var err error
			var clients *K8sClients

			clients, err = createK8sClients(args[0])
			if err != nil {
				return err
			}

			// ... get the map of contexts from the config file...
			if viper.Get("contexts") == nil {
				return fmt.Errorf("no k8s context specified in the config file: add at least one context")
			}

			contextsMap := viper.Get("contexts")

			// ... and deploy!
			err = deployInContext(args[0], clients, contextsMap.(map[string]interface{}))
			if err != nil {
				return err
			}

			return nil
		},
	}

	initCmd.Flags().StringVar(&opts.namespace, "namespace", "", "namespace")
	return initCmd
}

// deployInContext deploys the resources of the selected context in the config file
func deployInContext(context string, clients *K8sClients, contextsMap map[string]interface{}) error {
	// read the resources from the files assigned to the given context and call deployResources
	// for _, res := range contextsMap[context].([]interface{}) {
	//	...
	// }
	return nil
}

// deployResources ensures the existance of the namespace and calls the apply function for each resource
func deployResources(clients *K8sClients, namespace string, resources []Resource) error {
	// make sure the namespace exists
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
	// convert the resource GVK to GVR
	// gvr, err := fromGVKtoGVR(clients.discovery, res.Object.GroupVersionKind())
	// if err != nil {
	// 	return err
	// }

	// check if the resource exists on the cluster to create/patch it

	return nil
}
