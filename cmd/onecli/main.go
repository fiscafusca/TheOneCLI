package main

import (
	"os"

	"onecli/internal/onecli"
)

func main() {
	rootCmd := onecli.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
