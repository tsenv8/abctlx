package cmd

import (
	"abctlx/internal/airbyte"
	"context"

	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

var destCmd = &cobra.Command{
	Use:   "dest",
	Short: "Lists Destinations",
	Run: func(cmd *cobra.Command, args []string) {
		res := airbyte.NewAirbyteService(context.Background()).ListDestinations(nil)
		pretty.Print(res)
	},
}
var createDestFlags *airbyte.CreateDestinationFlags
var createDestCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates Destinations",
	Run: func(cmd *cobra.Command, args []string) {
		res := airbyte.NewAirbyteService(context.Background()).CreateDestination(*createDestFlags)
		pretty.Print(res)
	},
}
var updateDestFlags *airbyte.UpdateDestinationFlags
var updateDestCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates an existing Destination using its Destination Id",
	Run: func(cmd *cobra.Command, args []string) {
		res := airbyte.NewAirbyteService(context.Background()).UpdateDestination(*updateDestFlags)
		pretty.Print(res)
	},
}
var deleteDestName string
var deleteDestCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an existing Destination using its Destination Id",
	Run: func(cmd *cobra.Command, args []string) {
		res := airbyte.NewAirbyteService(context.Background()).DeleteDestination(deleteDestName)
		if res {
			pretty.Print("Deleted Successfully.")
		} else {
			pretty.Print("Failed Deletion of " + deleteDestName)
		}
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

func createDestCmdFlags() {
	createDestCmd.Flags().StringVar(&createDestFlags.Name, "name", "", "The Destination name.")
	createDestCmd.Flags().StringVar(&createDestFlags.ConfigType, "configType", "", "The Configuration Type.")

	// createDestCmd.Flags().String("srcName", "", "The name of the source to use.")
	// createDestCmd.Flags().String("destName", "", "The name of the destination to use.")
	// createDestCmd.Flags().String("schedule", "", "The schedule settings.")
	// createDestCmd.Flags().String("configType", "", "The configuration preset to use")
}

func updateDestCmdFlags() {
	updateDestCmd.Flags().StringVar(updateDestFlags.DestName, "destName", "", "The name of the destination to update.")
	updateDestCmd.Flags().StringVar(updateDestFlags.Name, "name", "", "Change destination name")
	updateDestCmd.Flags().StringVar(updateDestFlags.ConfigType, "configType", "", "Configuration type")

}

func deleteDestCmdFlags() {
	deleteDestCmd.Flags().StringVar(&deleteDestName, "name", "", "The destination name to delete")
}
