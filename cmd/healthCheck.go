package cmd

import (
	"abctlx/internal/airbyte"
	"abctlx/internal/config"
	"fmt"

	"github.com/spf13/cobra"
)

var healthCheckCmd = &cobra.Command{
	Use: "health",
	Run: func(cmd *cobra.Command, args []string) {
		output := airbyte.New(config.Data).HealthCheck()
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(healthCheckCmd)
}
