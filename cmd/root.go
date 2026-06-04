/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"abctlx/internal/config"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "abctlx",
	Short: "Airbyte Controls -- extended.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.LoadEnv)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
