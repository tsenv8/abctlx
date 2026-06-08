package cmd

import (
	"abctlx/internal/airbyte"
	"context"
	"fmt"

	"github.com/kr/pretty"
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
		updateParams := airbyte.UpdateSourceRequest{
			Configuration: &airbyte.UpdateSourceFields{},
		}

		CheckUpdateSourcesFlags(cmd, &updateParams)
		tsn, err := cmd.Flags().GetString("target-source")
		if err != nil {
			airbyte.NewAirbyteError("Field Required", "Target Source Name", err).Print()
		}

		res := airbyte.NewAirbyteService(context.Background()).UpdateSource(&updateParams, tsn)
		fmt.Println(res)
	},
}

var deleteSourceCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete source using source name.",
	Run: func(cmd *cobra.Command, args []string) {

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			airbyte.NewAirbyteError("Field Required", "Delete Source Name", err).Print()
		}

		success := airbyte.NewAirbyteService(context.Background()).DeleteSource(name)
		if success {
			pretty.Print("Deletion Successful")
		}
	},
}

func init() {
	rootCmd.AddCommand(sourcesCmd)
	sourcesCmd.AddCommand(createSourcesCmd)
	sourcesCmd.AddCommand(updateSourceCmd)
	sourcesCmd.AddCommand(deleteSourceCmd)
	createSourcesFlags()
	updateSourcesFlags()
	deleteSourcesFlags()
}

func deleteSourcesFlags() string {
	cmd := deleteSourceCmd
	var sourceName string
	cmd.Flags().String("name", "", "The target source")
	return sourceName
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

func CheckUpdateSourcesFlags(cmd *cobra.Command, updateParams *airbyte.UpdateSourceRequest) {
	conf := updateParams.Configuration
	errorField := "Update Sources Cmd"
	errorMsg := "Update Failed"

	if f := cmd.Flags().Lookup("name"); f != nil && f.Changed {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			airbyte.NewAirbyteError(errorMsg, errorField, err).Print()
		}

		updateParams.SourceName = name
	}

	if f := cmd.Flags().Lookup("db"); f != nil && f.Changed {
		dbName, err := cmd.Flags().GetString("db")
		if err != nil {
			airbyte.NewAirbyteError(errorMsg, errorField, err).Print()
		}

		conf.DBName = dbName
	}
	if f := cmd.Flags().Lookup("host"); f != nil && f.Changed {
		hostName, err := cmd.Flags().GetString("host")
		if err != nil {
			airbyte.NewAirbyteError(errorMsg, errorField, err).Print()
		}

		conf.HostName = hostName
	}

	if f := cmd.Flags().Lookup("pw"); f != nil && f.Changed {
		password, err := cmd.Flags().GetString("pw")
		if err != nil {
			airbyte.NewAirbyteError(errorMsg, errorField, err).Print()
		}

		conf.Password = password
	}
	if f := cmd.Flags().Lookup("pub"); f != nil && f.Changed {
		publicationName, err := cmd.Flags().GetString("pub")
		if err != nil {
			airbyte.NewAirbyteError(errorMsg, errorField, err).Print()

		}
		conf.ReplicationMethod.Publication = publicationName

	}
	if f := cmd.Flags().Lookup("rep"); f != nil && f.Changed {
		repSlotName, err := cmd.Flags().GetString("rep")
		if err != nil {
			airbyte.NewAirbyteError(errorMsg, errorField, err).Print()
		}

		conf.ReplicationMethod.ReplicationSlot = repSlotName
	}
	if f := cmd.Flags().Lookup("user"); f != nil && f.Changed {
		username, err := cmd.Flags().GetString("user")
		if err != nil {
			airbyte.NewAirbyteError(errorMsg, errorField, err).Print()
		}

		conf.Username = username
	}
	if f := cmd.Flags().Lookup("schema"); f != nil && f.Changed {
		schemas, _ := cmd.Flags().GetStringSlice("schema")
		if schemas == nil {
			schemas = []string{}
		}
		conf.Schemas = schemas
	}

	if f := cmd.Flags().Lookup("port"); f != nil && f.Changed {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			airbyte.NewAirbyteError(errorMsg, errorField, err).Print()
		}

		conf.Port = port
	}

	fmt.Printf("Updating only these fields: %+v\n", updateParams)
}

func updateSourcesFlags() {

	cmd := updateSourceCmd
	cmd.Flags().String("target-source", "", "Target Source Name")
	cmd.Flags().String("name", "", "Source Name")
	cmd.Flags().String("db", "", "Database Name")
	cmd.Flags().String("host", "", "Database Host Name")
	cmd.Flags().String("pw", "", "Database Password")
	cmd.Flags().String("pub", "", "Airbyte Publication Name")
	cmd.Flags().String("rep", "", "Airbyte Replication Slot Name")
	cmd.Flags().String("user", "", "Database Username")
	cmd.Flags().StringSlice("schema", []string{"public"}, "Database Schemas")
	cmd.Flags().Int("port", 2499, "Connection Port")
}
