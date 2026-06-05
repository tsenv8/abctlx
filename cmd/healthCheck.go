package cmd

import (
	"abctlx/internal/airbyte"
	"abctlx/internal/config"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var healthCheckCmd = &cobra.Command{
	Use: "health",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := airbyte.New(config.Data).HealthCheck()
		if err != nil {
			log.Fatal("ERROR: " + err.Error())
		}

		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(healthCheckCmd)
}
