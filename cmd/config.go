package cmd

import (
	"abctlx/internal/airbyte"
	"abctlx/internal/config"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Shows current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := airbyte.NewAirbyteClient(config.Data).GetConfig()
		apiEndpoint := airbyte.NewAirbyteClient(cfg).GetURL(nil)
		fmt.Println("-- Current Configuration --")
		fmt.Printf("\n URL:" + cfg.URL)
		fmt.Printf("\n API_URL:" + cfg.API_URL)
		fmt.Printf("\n API Endpoint:" + apiEndpoint)
		fmt.Printf("\n Port:" + strconv.Itoa(cfg.Port))
		fmt.Printf("\n ClientId:" + cfg.ClientId)
		fmt.Printf("\n ClientSecret:" + cfg.ClientKey)
		fmt.Println("\n---------------------------")

	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
