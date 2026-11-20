/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"abctlx/internal/airbyte"
	"abctlx/internal/config"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	PersistentPreRunE: skipProfile,
	Use:               "abctlx",
	Short:             "Airbyte Controls -- extended.",
}

var flagProfile string

// skipProfile resolves the chosen profile into config.Data before any command
// that talks to Airbyte. Profile management, help and completion do not need one.
func skipProfile(cmd *cobra.Command, args []string) error {
	top := cmd
	for top.Parent() != nil && top.Parent().Parent() != nil {
		top = top.Parent()
	}
	switch top.Name() {
	case "profile", "help", "completion", "abctlx":
		return nil
	}
	return config.Resolve(flagProfile)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagProfile, "profile", "P", "", "profile to use (default: the active one)")
	rootCmd.PersistentFlags().BoolVar(&airbyte.DryRun, "dry-run", false, "print every request that would change something, send none, and exit")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
