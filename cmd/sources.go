package cmd

import (
	"abctlx/helpers"
	"abctlx/internal/abctlx"
	"abctlx/internal/abctlx/cmd/sources"
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Interacts with sources",
	Run: func(cmd *cobra.Command, args []string) {
		runSources(cmd.Context())
	},
}

var parameters sources.CreateSourceParams
var createSourcesCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates source",
	Run: func(cmd *cobra.Command, args []string) {
		runCreateSources(cmd.Context())
	},
}

var updateSourceCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates source",
	Run: func(cmd *cobra.Command, args []string) {
		runUpdateSources(cmd.Context(), cmd)
	},
}

var deleteSourceCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete source using source name.",
	Run: func(cmd *cobra.Command, args []string) {
		runDeleteSources(cmd.Context(), cmd)
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

func runSources(ctx context.Context) {
	cmdHandler := abctlx.NewCmdHandler(ctx)
	res := cmdHandler.SrcSvc.ListSources(ctx, cmdHandler.AirbyteSvc.GetAccessToken())
	var body []table.Row

	for _, source := range res.Data {
		row := table.Row{
			source.SourceId,
			source.Name,
			source.SourceType,
			source.WorkspaceId,
		}

		body = append(body, row)
	}

	if len(body) < 1 {
		helpers.Log("info", "No sources found.")
	}

	helpers.Table(table.Row{"ID", "Name", "Type", "Workspace ID"}, body, "Sources")
}

func runCreateSources(ctx context.Context) {
	cmdHandler := abctlx.NewCmdHandler(ctx)
	res := cmdHandler.SrcSvc.CreateSource(
		ctx,
		parameters,
		*cmdHandler.AirbyteSvc.GetWorkspaceId(),
		cmdHandler.AirbyteSvc.GetAccessToken(),
	)

	helpers.Log("success", res.Name+" created.")
}

func runUpdateSources(ctx context.Context, cmd *cobra.Command) {
	updateParams := sources.UpdateSourceRequest{
		Configuration: &sources.UpdateSourceFields{},
	}

	CheckUpdateSourcesFlags(cmd, &updateParams)
	tsn, err := cmd.Flags().GetString("target-source")
	if err != nil {
		helpers.Error("Target Source Name Field Required", err)
	}

	cmdHandler := abctlx.NewCmdHandler(ctx)
	res := cmdHandler.SrcSvc.UpdateSource(
		ctx,
		&updateParams,
		tsn,
		cmdHandler.AirbyteSvc.GetAccessToken(),
	)

	helpers.Log("success", res.Name+" updated.")
}

func runDeleteSources(ctx context.Context, cmd *cobra.Command) {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		helpers.Error("Delete Source Name Field Required", err)
	}

	cmdHandler := abctlx.NewCmdHandler(ctx)
	res := cmdHandler.SrcSvc.DeleteSource(
		ctx,
		name,
		cmdHandler.AirbyteSvc.GetAccessToken(),
	)

	helpers.BoolLog(res)
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

func CheckUpdateSourcesFlags(cmd *cobra.Command, updateParams *sources.UpdateSourceRequest) {
	conf := updateParams.Configuration
	errorField := "Update Sources Cmd"
	errorMsg := "Update Failed"

	if f := cmd.Flags().Lookup("name"); f != nil && f.Changed {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			helpers.Error(errorMsg+" "+errorField, err)
			// airbyte.NewAirbyteError(errorMsg, errorField, err).Print()
		}

		updateParams.SourceName = name
	}

	if f := cmd.Flags().Lookup("db"); f != nil && f.Changed {
		dbName, err := cmd.Flags().GetString("db")
		if err != nil {
			helpers.Error(errorMsg+" "+errorField, err)
		}

		conf.DBName = dbName
	}
	if f := cmd.Flags().Lookup("host"); f != nil && f.Changed {
		hostName, err := cmd.Flags().GetString("host")
		if err != nil {
			helpers.Error(errorMsg+" "+errorField, err)

		}

		conf.HostName = hostName
	}

	if f := cmd.Flags().Lookup("pw"); f != nil && f.Changed {
		password, err := cmd.Flags().GetString("pw")
		if err != nil {
			helpers.Error(errorMsg+" "+errorField, err)
		}

		conf.Password = password
	}
	if f := cmd.Flags().Lookup("pub"); f != nil && f.Changed {
		publicationName, err := cmd.Flags().GetString("pub")
		if err != nil {
			helpers.Error(errorMsg+" "+errorField, err)
		}
		conf.ReplicationMethod.Publication = publicationName

	}
	if f := cmd.Flags().Lookup("rep"); f != nil && f.Changed {
		repSlotName, err := cmd.Flags().GetString("rep")
		if err != nil {
			helpers.Error(errorMsg+" "+errorField, err)
		}

		conf.ReplicationMethod.ReplicationSlot = repSlotName
	}
	if f := cmd.Flags().Lookup("user"); f != nil && f.Changed {
		username, err := cmd.Flags().GetString("user")
		if err != nil {
			helpers.Error(errorMsg+" "+errorField, err)
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
			helpers.Error(errorMsg+" "+errorField, err)
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
