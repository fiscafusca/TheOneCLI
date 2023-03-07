package onecli

import (
	"github.com/spf13/cobra"
)

type options struct {
	name        string
	namespace   string
	configFile  string
	projectPath string
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

	rootCmd.PersistentFlags().StringVar(&opts.configFile, "config", "", "custom config file")

	rootCmd.AddCommand(NewInitCommand())
	rootCmd.AddCommand(NewDeployCommand())

	return rootCmd
}
