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
		foo()
	},
}

func foo() {
	s := airbyte.NewAirbyteService(context.Background())
	var body []table.Row
	res := s.ListConnections(nil)
	for _, data := range res.Data {
		rowString := []string{data.ConnectionId, data.Name, data.SourceId, data.DestinationId}
		row := abctlx.ToRow(rowString)
		body = append(body, row)
	}
	header := table.Row{"ID", "Name", "Source ID", "Destination ID"}
	abctlx.Table(header, body)
}

func init() {
	rootCmd.AddCommand(conCmd)
}
