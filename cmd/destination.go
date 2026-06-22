package cmd

import (
	"abctlx/helpers"
	"abctlx/internal/abctlx"
	"abctlx/internal/abctlx/cmd/destination"
	"context"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var destCmd = &cobra.Command{
	Use:   "dest",
	Short: "Lists Destinations",
	Run: func(cmd *cobra.Command, args []string) {
		runDest(cmd.Context())
	},
}

var createDestFlags destination.CreateDestinationFlags
var createDestCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates Destinations",
	Run: func(cmd *cobra.Command, args []string) {
		runCreateDest(cmd.Context())
	},
}

var updateDestFlags destination.UpdateDestinationFlags
var updateDestCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates an existing Destination using its Destination Id",
	Run: func(cmd *cobra.Command, args []string) {
		runUpdateDest(cmd.Context())
	},
}

var deleteDestName string
var deleteDestCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an existing Destination using its Destination Id",
	Run: func(cmd *cobra.Command, args []string) {
		runDeleteDest(cmd.Context())
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

func runDest(ctx context.Context) {

	cmdHandler := abctlx.NewCmdHandler(ctx)
	res := cmdHandler.DestSvc.ListDestinations(
		ctx,
		nil,
		cmdHandler.AirbyteSvc.GetAccessToken(),
	)

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

	helpers.Table(table.Row{"ID", "Name", "Type", "Workspace ID"}, body, "Destinations")
}

func runCreateDest(ctx context.Context) {
	cmdHandler := abctlx.NewCmdHandler(ctx)
	res := cmdHandler.DestSvc.CreateDestination(
		ctx,
		createDestFlags,
		cmdHandler.AirbyteSvc.GetAccessToken(),
		*cmdHandler.AirbyteSvc.GetWorkspaceId(),
	)

	helpers.Log("success", res.Name+" created.")
}

func runUpdateDest(ctx context.Context) {
	// abSvc := airbyte.NewAirbyteService(ctx)
	// destSvc := destination.NewService()
	cmdHandler := abctlx.NewCmdHandler(ctx)
	res := cmdHandler.DestSvc.UpdateDestination(
		ctx,
		updateDestFlags,
		cmdHandler.AirbyteSvc.GetAccessToken(),
	)
	// res := destSvc.UpdateDestination(ctx, updateDestFlags, abSvc.GetAccessToken())
	// res := airbyte.NewAirbyteService(context.Background()).UpdateDestination(updateDestFlags)
	// pretty.Print(res)

	helpers.Log("success", res.Name+" updated.")
}

func runDeleteDest(ctx context.Context) {
	cmdHandler := abctlx.NewCmdHandler(ctx)
	res := cmdHandler.DestSvc.DeleteDestination(
		ctx,
		deleteDestName,
		cmdHandler.AirbyteSvc.GetAccessToken(),
	)
	helpers.BoolLog(res)
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
