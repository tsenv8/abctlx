package cmd

import (
	"abctlx/internal/airbyte"
	"abctlx/internal/config"
	"fmt"

	"github.com/spf13/cobra"
)

var generateAccessTokenCmd = &cobra.Command{
	Use: "token",
	Run: func(cmd *cobra.Command, args []string) {
		output := airbyte.New(config.Data).GenerateAccessToken()
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(generateAccessTokenCmd)
}
