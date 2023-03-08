package onecli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type options struct {
	name      string
	namespace string
}

var opts = &options{}

func NewRootCommand() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "onecli",
		Short: "One CLI to rule them all",
		Long: `One CLI to rule them all,
One CLI to find them,
One CLI to bring them all and in the cluster deploy them.`,
	}

	rootCmd.AddCommand(NewInitCommand())
	rootCmd.AddCommand(NewDeployCommand())
	rootCmd.AddCommand(NewAddContextCommand())

	return rootCmd
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("onecli")
	viper.AddConfigPath(currentPath)

	contextMap := make(map[string]interface{})
	viper.Set("k8s-contexts", contextMap)

	if err := viper.SafeWriteConfig(); err != nil {
		viper.SafeWriteConfigAs(currentPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error loading file:", viper.ConfigFileUsed())
		os.Exit(1)
	}
}
