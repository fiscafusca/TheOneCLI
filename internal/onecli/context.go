package onecli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewAddContextCommand() *cobra.Command {
	initCmd := &cobra.Command{

		Use:   "add-context [NAME]",
		Short: "Add a new context to the config file",
		Long: `A veeeeeeeeeeeeeeeeeeeeery
loooooooooooooooooooooooooooooooooooong
descriptiooooooooooooooooooooooooooooon.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			var contextsMap map[string]interface{}

			newContextName := args[0]

			fmt.Printf("Adding new context: %s ...\n", newContextName)

			err := viper.ReadInConfig()
			if err != nil {
				return fmt.Errorf("error reading config file: %w", err)
			}

			if viper.Get("contexts") == nil {
				contextsMap = make(map[string]interface{})
			} else {
				contextsMap = viper.Get("contexts").(map[string]interface{})
			}

			if contextsMap[newContextName] == nil {
				contextsMap[newContextName] = []string{}
			}

			fmt.Println(contextsMap)

			viper.Set("contexts", contextsMap)

			return viper.WriteConfig()
		},
	}

	return initCmd
}
