package cmd

import (
	"abctlx/internal/abctlx"
	"abctlx/internal/airbyte"
	"context"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var conCmd = &cobra.Command{
	Use:   "connection",
	Short: "Retrieves all Airbyte Connections",
	Run: func(cmd *cobra.Command, args []string) {
		s := airbyte.NewAirbyteService(context.Background())
		var body []table.Row
		res := s.ListConnections(nil)
		for _, data := range res.Data {
			rowString := []string{data.ConnectionId, data.Name, data.SourceId, data.DestinationId}
			row := abctlx.ToRow(rowString)
			body = append(body, row)
		}
		header := table.Row{"ID", "Name", "Source ID", "Destination ID"}
		abctlx.Table(header, body, "Connections")
	},
}

var createConReq airbyte.CreateConnectionRequest
var createConCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a connection between source and destination.",
	Run: func(cmd *cobra.Command, args []string) {
		airbyte.NewAirbyteService(context.Background()).CreateConnection(&createConReq)

	},
}

func init() {
	rootCmd.AddCommand(conCmd)
	conCmd.AddCommand(createConCmd)
	createConCmdFlags()
}

func createConCmdFlags() {
	createConCmd.Flags().StringVar(&createConReq.Body.Name, "name", "", "The name of the Connection.")
	createConCmd.Flags().StringVar(&createConReq.Body.SourceId, "source", "", "The Source ID to connect.")
	createConCmd.Flags().StringVar(&createConReq.Body.DestinationId, "dest", "", "The Destination ID to connect.")
}
