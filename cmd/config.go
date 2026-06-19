package cmd

import (
	"abctlx/internal/abctlx"
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
		headers := abctlx.ToRow([]string{"Variable", "Value"})
		body := []table.Row{
			abctlx.ToRow([]string{"URL", cfg.URL}),
			abctlx.ToRow([]string{"API URL", cfg.API_URL}),
			abctlx.ToRow([]string{"API Endpoint", apiEndpoint}),
			abctlx.ToRow([]string{"Port", strconv.Itoa(cfg.Port)}),
			abctlx.ToRow([]string{"Client ID", cfg.ClientId}),
			abctlx.ToRow([]string{"Client Secret", cfg.ClientKey}),
			abctlx.ToRow([]string{"Kubeconfig", cfg.Kubeconfig}),
			abctlx.ToRow([]string{"Namespace", cfg.Namespace}),
		}
		abctlx.Table(headers, body, "Configuration")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
