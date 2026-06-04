package cmd

import (
	"abctlx/internal/airbyte"
	"abctlx/internal/config"
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use: "config",
	Run: func(cmd *cobra.Command, args []string) {
		output := airbyte.New(config.Data).GetURL(nil)
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
