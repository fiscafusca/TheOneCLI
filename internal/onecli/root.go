package onecli

import (
	"github.com/spf13/cobra"
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

	// add the cobra commands

	return rootCmd
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {

	// initialize the viper config file in the current directory

}
