package cmd

import (
	"abctlx/helpers"
	"abctlx/internal/airbyte"
	"abctlx/internal/config"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Shows current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := airbyte.NewAirbyteClient(config.Data).GetConfig()
		apiEndpoint := airbyte.NewAirbyteClient(cfg).GetURL(nil)
		headers := helpers.ToRow([]string{"Variable", "Value"})
		body := []table.Row{
			helpers.ToRow([]string{"URL", cfg.URL}),
			helpers.ToRow([]string{"API URL", cfg.API_URL}),
			helpers.ToRow([]string{"API Endpoint", apiEndpoint}),
			helpers.ToRow([]string{"Port", strconv.Itoa(cfg.Port)}),
			helpers.ToRow([]string{"Client ID", cfg.ClientId}),
			helpers.ToRow([]string{"Client Secret", cfg.ClientKey}),
			helpers.ToRow([]string{"Kubeconfig", cfg.Kubeconfig}),
			helpers.ToRow([]string{"Namespace", cfg.Namespace}),
		}
		helpers.Table(headers, body, "Configuration")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
