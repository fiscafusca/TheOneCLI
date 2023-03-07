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
	opts.projectPath = path.Join(currentPath, opts.name)
	if err := os.Mkdir(opts.projectPath, fs.ModePerm); err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}
	resourcesDir := path.Join(opts.projectPath, "resources")
	if err := os.Mkdir(resourcesDir, fs.ModePerm); err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}

	viper.AddConfigPath(opts.projectPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("onecli")

	contextMap := make(map[string]interface{})
	contextMap["default"] = []string{}
	viper.Set("k8s-contexts", contextMap)

	if err := viper.SafeWriteConfig(); err != nil {
		viper.SafeWriteConfigAs(opts.projectPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error loading file:", viper.ConfigFileUsed())
	}
	return nil
}
