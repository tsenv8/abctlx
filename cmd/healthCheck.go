package cmd

import (
	"abctlx/internal/airbyte"
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var healthCheckCmd = &cobra.Command{
	Use: "health",
	Run: func(cmd *cobra.Command, args []string) {
		var status string
		res, err := airbyte.NewAirbyteService(context.Background()).Health()
		if err != nil {
			log.Fatal("ERROR:" + err.Error())
		}

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
