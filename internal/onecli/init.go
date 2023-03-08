package onecli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	initCmd.Flags().StringVar(&opts.name, "name", defaultName, "project name")
	viper.BindPFlag("name", initCmd.Flags().Lookup("name"))
	return initCmd
}

func initProject() error {
	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	resourcesDir := path.Join(currentPath, "resources")
	if err := os.Mkdir(resourcesDir, fs.ModePerm); err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}
	return nil
}
