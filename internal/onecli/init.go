package onecli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const defaultName = "onecli"

func NewInitCommand() *cobra.Command {
	initCmd := &cobra.Command{

		Use:   "init",
		Short: "Initialize the project",
		Long: `A veeeeeeeeeeeeeeeeeeeeery
loooooooooooooooooooooooooooooooooooong
descriptiooooooooooooooooooooooooooooon.`,

		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Println("Initializing...")
			return initProject()
		},
	}

	// add and bind the name flag

	return initCmd
}

func initProject() error {

	// initialize the resources folder inside the current directory

	return nil
}
