package cmd

import (
	"abctlx/internal/abctlx"
	"abctlx/internal/airbyte"
	"context"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

var destCmd = &cobra.Command{
	Use:   "dest",
	Short: "Lists Destinations",
	Run: func(cmd *cobra.Command, args []string) {
		runDest()
	},
}

var createDestFlags airbyte.CreateDestinationFlags
var createDestCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates Destinations",
	Run: func(cmd *cobra.Command, args []string) {
		runCreateDest()
	},
}

var updateDestFlags airbyte.UpdateDestinationFlags
var updateDestCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates an existing Destination using its Destination Id",
	Run: func(cmd *cobra.Command, args []string) {
		runUpdateDest()
	},
}

var deleteDestName string
var deleteDestCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an existing Destination using its Destination Id",
	Run: func(cmd *cobra.Command, args []string) {
		runDeleteDest()
	},
}

func init() {
	rootCmd.AddCommand(destCmd)

	createDestCmdFlags()
	deleteDestCmdFlags()
	updateDestCmdFlags()

	destCmd.AddCommand(createDestCmd)
	destCmd.AddCommand(deleteDestCmd)
	destCmd.AddCommand(updateDestCmd)
}

func runDest() {
	res := airbyte.NewAirbyteService(context.Background()).ListDestinations(nil)
	var body []table.Row
	for _, data := range res.Data {
		row := table.Row{
			data.DestinationId,
			data.Name,
			data.DestinationType,
			data.WorkspaceId,
		}

		body = append(body, row)
	}

	abctlx.Table(table.Row{"ID", "Name", "Type", "Workspace ID"}, body, "Destinations")
}

func runCreateDest() {
	res := airbyte.NewAirbyteService(context.Background()).CreateDestination(createDestFlags)
	pretty.Print(res)
}

func runUpdateDest() {
	res := airbyte.NewAirbyteService(context.Background()).UpdateDestination(updateDestFlags)
	pretty.Print(res)
}

func runDeleteDest() {
	res := airbyte.NewAirbyteService(context.Background()).DeleteDestination(deleteDestName)
	abctlx.BoolLog(res)
}

func createDestCmdFlags() {
	createDestCmd.Flags().StringVar(&createDestFlags.Name, "name", "", "The Destination name.")
	createDestCmd.Flags().StringVar(&createDestFlags.ConfigType, "configType", "", "The Configuration Type.")

	// createDestCmd.Flags().String("srcName", "", "The name of the source to use.")
	// createDestCmd.Flags().String("destName", "", "The name of the destination to use.")
	// createDestCmd.Flags().String("schedule", "", "The schedule settings.")
	// createDestCmd.Flags().String("configType", "", "The configuration preset to use")
}

func updateDestCmdFlags() {
	updateDestCmd.Flags().StringVar(&updateDestFlags.DestName, "destName", "", "The name of the destination to update.")
	updateDestCmd.Flags().StringVar(&updateDestFlags.Name, "name", "", "Change destination name")
	updateDestCmd.Flags().StringVar(&updateDestFlags.ConfigType, "configType", "", "Configuration type")
}

func deleteDestCmdFlags() {
	deleteDestCmd.Flags().StringVar(&deleteDestName, "name", "", "The destination name to delete")
}
