package cmd

import (
	"abctlx/internal/airbyte"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var healthCheckCmd = &cobra.Command{
	Use:   "health",
	Short: "Health Check for the API",
	Run: func(cmd *cobra.Command, args []string) {
		var status string
		res := airbyte.NewAirbyteService(context.Background()).Health()

		if res.Status {
			status = "true"
		} else {
			status = "false"
		}

		fmt.Println("Status:" + status)
	},
}

func init() {
	rootCmd.AddCommand(healthCheckCmd)
}
