package cmd

import (
	"abctlx/internal/airbyte"
	"abctlx/internal/config"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var generateAccessTokenCmd = &cobra.Command{
	Use: "token",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := airbyte.New(config.Data).GenerateAccessToken()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(generateAccessTokenCmd)
}
