package onecli

import (
	"github.com/spf13/cobra"
)

type FlagPole struct {
	Name   string
	Config string
}

// NewRootCommand returns a new cobra.Command for root command
func NewRootCommand() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "onecli",
		Short: "A CLI to rule them all",
	}

	return rootCmd
}
