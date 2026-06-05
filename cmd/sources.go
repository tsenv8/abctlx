package cmd

import (
	"abctlx/internal/airbyte"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var parameters airbyte.CreateSourceParams

var sourcesCmd = &cobra.Command{
	Use: "sources",
	Run: func(cmd *cobra.Command, args []string) {
		res := airbyte.NewAirbyteService(context.Background()).ListSources()
		fmt.Println(res.Data)
	},
}

var createSourcesCmd = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(parameters)
		// airbyte.NewAirbyteService(context.Background()).CreateSource(parameters)
	},
}

func init() {
	//create
	rootCmd.AddCommand(sourcesCmd)
	sourcesCmd.AddCommand(createSourcesCmd)
	createSourcesCmd.Flags().StringVar(&parameters.Name, "name", "sourcedb", "Source Name")
	createSourcesCmd.Flags().StringVar(&parameters.DBName, "db", "postgres", "Database Name")
	createSourcesCmd.Flags().StringVar(&parameters.HostName, "host", "postgres", "Database Host Name")
	createSourcesCmd.Flags().StringVar(&parameters.Password, "pw", "1", "Database Password")
	createSourcesCmd.Flags().StringVar(&parameters.PublicationName, "pub", "airbyte_publication", "Airbyte Publication Name")
	createSourcesCmd.Flags().StringVar(&parameters.ReplicationSlot, "rep", "airbyte_slot", "Airbyte Replication Slot Name")
	createSourcesCmd.Flags().StringVar(&parameters.Username, "user", "postgres", "Database Username")
	// createSourcesCmd.Flags().StringVar(&parameters.Schemas, "schema", ["public"] ,  "Database Username")
	createSourcesCmd.Flags().StringSliceVar(&parameters.Schemas, "schema", []string{"public"}, "Database Schemas default: public")
	createSourcesCmd.Flags().IntVar(&parameters.Port, "port", 2499, "Connection Port")

}
