package onecli

import (
	"github.com/spf13/cobra"
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

			// add new k8s contexts in the viper config file

			return nil
		},
	}

	return initCmd
}
