package cmd

import (
	"abctlx/internal/airbyte"
	"abctlx/internal/config"
	"fmt"

	"github.com/spf13/cobra"
)

// generateAccessTokenCmd represents the generateAccessToken command
var generateAccessTokenCmd = &cobra.Command{
	Use: "generateAccessToken",
	Run: func(cmd *cobra.Command, args []string) { run() },
}

func init() {
	rootCmd.AddCommand(generateAccessTokenCmd)
}

func run() {
	token := airbyte.New(config.Data).GenerateAccessToken()
	fmt.Println(token)
}
