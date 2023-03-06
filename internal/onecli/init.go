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

var name string

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

	initCmd.Flags().StringVar(&flags.name, "name", defaultName, "project name")
	viper.BindPFlag("name", initCmd.Flags().Lookup("name"))
	return initCmd
}

func initProject() error {
	if flags.configFile != "" {
		viper.SetConfigFile(flags.configFile)
	} else {

		currentPath, err := os.Getwd()
		if err != nil {
			return err
		}
		projectPath := path.Join(currentPath, name)
		if err := os.Mkdir(projectPath, fs.ModePerm); err != nil && !errors.Is(err, fs.ErrExist) {
			return err
		}

		viper.AddConfigPath(projectPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("onecli")

		if err := viper.SafeWriteConfig(); err != nil {
			viper.SafeWriteConfigAs(projectPath)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error loading file:", viper.ConfigFileUsed())
	}
	return nil
}
