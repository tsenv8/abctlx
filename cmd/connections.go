package cmd

import (
	"abctlx/helpers"
	"abctlx/internal/airbyte"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

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
			var tags []string
			for _, t := range data.Tags {
				tags = append(tags, t.Name)
			}
			rowString := []string{data.ConnectionId, data.Name, strings.Join(tags, ","), data.Status, data.SourceId, data.DestinationId}
			row := helpers.ToRow(rowString)
			body = append(body, row)
		}
		header := table.Row{"ID", "Name", "Tags", "Status", "Source ID", "Destination ID"}
		helpers.Table(header, body, "Connections")
	},
}

var createConReq airbyte.CreateConnectionRequest
var createConCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a connection between source and destination.",
	Run: func(cmd *cobra.Command, args []string) {
		airbyte.NewAirbyteService(context.Background()).CreateConnection(createConReq)

	},
}

var streamsConId string
var streamsAvailable bool
var streamsRaw bool
var streamsRefresh bool
var streamsCmd = &cobra.Command{
	Use:   "streams",
	Short: "Read-only: lists the enabled streams of a connection (--available: what the source offers).",
	Run: func(cmd *cobra.Command, args []string) {
		if streamsConId == "" {
			helpers.Error("--id is required", nil)
		}
		s := airbyte.NewAirbyteService(context.Background())
		con, conRaw, err := s.GetConnectionByID(streamsConId)
		if err != nil {
			helpers.Error("Get Connection Failed", err)
		}

		if !streamsAvailable {
			if streamsRaw {
				fmt.Println(string(conRaw))
				return
			}
			var body []table.Row
			for _, raw := range con.Configurations.Streams {
				var st airbyte.StreamConfiguration
				if err := json.Unmarshal(raw, &st); err != nil {
					helpers.Error("Stream JSON Unmarshal Failed", err)
				}
				body = append(body, table.Row{st.Namespace, st.Name, st.SyncMode, strings.Join(st.CursorField, "."), fmt.Sprint(st.PrimaryKey)})
			}
			sort.Slice(body, func(i, j int) bool {
				return fmt.Sprint(body[i][0], ".", body[i][1]) < fmt.Sprint(body[j][0], ".", body[j][1])
			})
			helpers.Table(table.Row{"Namespace", "Name", "Sync mode", "Cursor", "Primary key"}, body, fmt.Sprintf("Enabled streams: %s (%d)", con.Name, len(body)))
			return
		}

		avail, availRaw, err := s.ListAvailableStreams(con.SourceId, con.DestinationId, streamsRefresh)
		if err != nil {
			helpers.Error("List Available Streams Failed", err)
		}
		if streamsRaw {
			fmt.Println(string(availRaw))
			return
		}
		var body []table.Row
		for _, a := range avail {
			body = append(body, table.Row{a.StreamNamespace, a.StreamName, strings.Join(a.SyncModes, ","), fmt.Sprint(a.SourceDefinedPrimaryKey)})
		}
		helpers.Table(table.Row{"Namespace", "Name", "Sync modes", "Source PK"}, body, fmt.Sprintf("Available streams: %s (%d)", con.Name, len(body)))
	},
}

func init() {
	rootCmd.AddCommand(conCmd)
	conCmd.AddCommand(createConCmd)
	conCmd.AddCommand(streamsCmd)
	streamsCmd.Flags().StringVar(&streamsConId, "id", "", "Connection ID.")
	streamsCmd.Flags().BoolVar(&streamsAvailable, "available", false, "List the streams the source offers instead of the enabled ones.")
	streamsCmd.Flags().BoolVar(&streamsRefresh, "refresh", false, "With --available: ask Airbyte to rediscover the source schema instead of using its cache.")
	streamsCmd.Flags().BoolVar(&streamsRaw, "json", false, "Print the raw API response (used to check the real shape).")
	createConCmdFlags()
}

func createConCmdFlags() {
	createConCmd.Flags().StringVar(&createConReq.Body.Name, "name", "", "The name of the Connection.")
	createConCmd.Flags().StringVar(&createConReq.Body.SourceId, "source", "", "The Source ID to connect.")
	createConCmd.Flags().StringVar(&createConReq.Body.DestinationId, "dest", "", "The Destination ID to connect.")
}
