package main

import "github.com/spf13/cobra"

// DefaultNodeHome defines the daemon's home directory used if none is provided.
const DefaultNodeHome = ".usdWz"

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "usdwzd",
		Short: "usdWz blockchain daemon",
	}

	return rootCmd
}

func execute() error {
	return newRootCmd().Execute()
}
