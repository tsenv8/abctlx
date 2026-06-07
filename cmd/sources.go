package cmd

import (
	"abctlx/internal/airbyte"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var parameters airbyte.CreateSourceParams

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Interacts with sources",
	Run: func(cmd *cobra.Command, args []string) {
		res := airbyte.NewAirbyteService(context.Background()).ListSources()
		fmt.Println(res.Data)
	},
}

var createSourcesCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates source",
	Run: func(cmd *cobra.Command, args []string) {
		airbyte.NewAirbyteService(context.Background()).CreateSource(parameters)
	},
}

var updateSourceCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates source",
	Run: func(cmd *cobra.Command, args []string) {
		updateParams := airbyte.CreateSourceParams{}
		CheckUpdateSourcesFlags(cmd, &updateParams)
		// airbyte.NewAirbyteService(context.Background()).CreateSource(parameters)
	},
}

func init() {
	//create
	rootCmd.AddCommand(sourcesCmd)
	sourcesCmd.AddCommand(createSourcesCmd)
	sourcesCmd.AddCommand(updateSourceCmd)
	createSourcesFlags()
	updateSourcesFlags()
}

func createSourcesFlags() {
	cmd := createSourcesCmd
	cmd.Flags().StringVar(&parameters.Name, "name", "sourcedb", "Source Name")
	cmd.Flags().StringVar(&parameters.DBName, "db", "postgres", "Database Name")
	cmd.Flags().StringVar(&parameters.HostName, "host", "localhost", "Database Host Name")
	cmd.Flags().StringVar(&parameters.Password, "pw", "1", "Database Password")
	cmd.Flags().StringVar(&parameters.PublicationName, "pub", "airbyte_publication", "Airbyte Publication Name")
	cmd.Flags().StringVar(&parameters.ReplicationSlot, "rep", "airbyte_slot", "Airbyte Replication Slot Name")
	cmd.Flags().StringVar(&parameters.Username, "user", "postgres", "Database Username")
	cmd.Flags().StringSliceVar(&parameters.Schemas, "schema", []string{"public"}, "Database Schemas")
	cmd.Flags().IntVar(&parameters.Port, "port", 2499, "Connection Port")
}

func CheckUpdateSourcesFlags(cmd *cobra.Command, updateParams *airbyte.CreateSourceParams) {

	// 2. Manually check which flags were actually passed by the user
	if f := cmd.Flags().Lookup("name"); f != nil && f.Changed {
		updateParams.Name, _ = cmd.Flags().GetString("name")
	}
	if f := cmd.Flags().Lookup("db"); f != nil && f.Changed {
		updateParams.DBName, _ = cmd.Flags().GetString("db")
	}
	if f := cmd.Flags().Lookup("host"); f != nil && f.Changed {
		updateParams.HostName, _ = cmd.Flags().GetString("host")
	}
	if f := cmd.Flags().Lookup("pw"); f != nil && f.Changed {
		updateParams.Password, _ = cmd.Flags().GetString("pw")
	}
	if f := cmd.Flags().Lookup("pub"); f != nil && f.Changed {
		updateParams.PublicationName, _ = cmd.Flags().GetString("pub")
	}
	if f := cmd.Flags().Lookup("rep"); f != nil && f.Changed {
		updateParams.ReplicationSlot, _ = cmd.Flags().GetString("rep")
	}
	if f := cmd.Flags().Lookup("user"); f != nil && f.Changed {
		updateParams.Username, _ = cmd.Flags().GetString("user")
	}
	if f := cmd.Flags().Lookup("schema"); f != nil && f.Changed {
		updateParams.Schemas, _ = cmd.Flags().GetStringSlice("schema")
	}
	if f := cmd.Flags().Lookup("port"); f != nil && f.Changed {
		updateParams.Port, _ = cmd.Flags().GetInt("port")
	}

	fmt.Printf("Updating only these fields: %+v\n", updateParams)
	// airbyte.NewAirbyteService(context.Background()).UpdateSource(updateParams)
}

func updateSourcesFlags() {

	cmd := updateSourceCmd
	cmd.Flags().String("name", "", "Source Name")
	cmd.Flags().String("db", "", "Database Name")
	cmd.Flags().String("host", "", "Database Host Name")
	cmd.Flags().String("pw", "", "Database Password")
	cmd.Flags().String("pub", "", "Airbyte Publication Name")
	cmd.Flags().String("rep", "", "Airbyte Replication Slot Name")
	cmd.Flags().String("user", "", "Database Username")
	cmd.Flags().StringSlice("schema", []string{}, "Database Schemas")
	cmd.Flags().Int("port", 2499, "Connection Port")

	// cmd := updateSourceCmd
	// cmd.Flags().StringVar(&parameters.Name, "name", "" , "Source Name")
	// cmd.Flags().StringVar(&parameters.DBName, "db", "", "Database Name")
	// cmd.Flags().StringVar(&parameters.HostName, "host", "", "Database Host Name")
	// cmd.Flags().StringVar(&parameters.Password, "pw", "", "Database Password")
	// cmd.Flags().StringVar(&parameters.PublicationName, "pub", "", "Airbyte Publication Name")
	// cmd.Flags().StringVar(&parameters.ReplicationSlot, "rep", "", "Airbyte Replication Slot Name")
	// cmd.Flags().StringVar(&parameters.Username, "user", "", "Database Username")
	// cmd.Flags().StringSliceVar(&parameters.Schemas, "schema", []string{}, "Database Schemas")
	// cmd.Flags().IntVar(&parameters.Port, "port", 0, "Connection Port")
}
