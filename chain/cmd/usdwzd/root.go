package main

import "github.com/spf13/cobra"

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
