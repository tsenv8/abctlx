package cmd

import (
	"abctlx/helpers"
	"abctlx/internal/airbyte"
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	addStreamsId   string
	addDryRun      bool
	addSkipMissing bool
	addYes         bool
	addNoRefresh   bool
)

var addStreamsCmd = &cobra.Command{
	Use:   "add-streams <schema.table>...",
	Short: "Enable more tables (streams) on a connection, keeping the existing ones",
	Long: `Reads the connection's current streams, adds the given schema.table names with a sync
mode chosen per table from what the source reports, and sends the full list back (Airbyte
replaces the whole list on update, so the existing streams are sent back unchanged).
Schedule and every other setting are not touched.

The source schema is rediscovered first (use --no-refresh to use Airbyte's cache). Tables
already enabled are skipped. A table the source does not offer stops the command, unless
--skip-missing. It prints the plan, then asks before applying; --dry-run stops after the plan.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if addStreamsId == "" {
			return fmt.Errorf("--id is required")
		}
		s := airbyte.NewAirbyteService(context.Background())
		con, _, err := s.GetConnectionByID(addStreamsId)
		if err != nil {
			return fmt.Errorf("get connection: %w", err)
		}
		avail, _, err := s.ListAvailableStreams(con.SourceId, con.DestinationId, !addNoRefresh)
		if err != nil {
			return fmt.Errorf("list available streams: %w", err)
		}
		plan, streams, err := airbyte.PlanAddStreams(con.Configurations.Streams, avail, args)
		if err != nil {
			return err
		}

		var body []table.Row
		for _, p := range plan.Add {
			body = append(body, table.Row{p.Table, p.SyncMode, strings.Join(p.Cursor, "."), fmt.Sprint(p.PrimaryKey), p.Reason})
		}
		helpers.Table(table.Row{"Table", "Sync mode", "Cursor", "Primary key", "Why"}, body,
			fmt.Sprintf("Add to %s: %d now enabled -> %d after", con.Name, len(con.Configurations.Streams), len(streams)))
		if len(plan.Already) > 0 {
			fmt.Printf("\nAlready enabled (skipped): %s\n", strings.Join(plan.Already, ", "))
		}
		if len(plan.Missing) > 0 {
			fmt.Printf("\nNot offered by the source: %s\n", strings.Join(plan.Missing, ", "))
			if !addSkipMissing {
				return fmt.Errorf("nothing changed: %d table(s) are missing (publish them first, or pass --skip-missing)", len(plan.Missing))
			}
		}
		if len(plan.Add) == 0 {
			fmt.Println("\nNothing to add.")
			return nil
		}
		if addDryRun {
			fmt.Println("\nDry run: nothing sent.")
			return nil
		}
		if !addYes {
			fmt.Printf("\nApply to connection %s (%s)? [y/N] ", con.Name, con.ConnectionId)
			ans, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			if a := strings.ToLower(strings.TrimSpace(ans)); a != "y" && a != "yes" {
				fmt.Println("Cancelled; nothing changed.")
				return nil
			}
		}
		updated, err := s.UpdateConnectionStreams(addStreamsId, streams)
		if err != nil {
			return fmt.Errorf("update connection: %w", err)
		}
		fmt.Printf("\nDone: connection now has %d streams.\n", len(updated.Configurations.Streams))
		if len(updated.Configurations.Streams) != len(streams) {
			return fmt.Errorf("expected %d streams after the update, Airbyte reports %d: check the connection in the UI", len(streams), len(updated.Configurations.Streams))
		}
		return nil
	},
}

func init() {
	conCmd.AddCommand(addStreamsCmd)
	f := addStreamsCmd.Flags()
	f.StringVar(&addStreamsId, "id", "", "Connection ID.")
	f.BoolVar(&addDryRun, "dry-run", false, "Print the plan and stop.")
	f.BoolVar(&addSkipMissing, "skip-missing", false, "Skip tables the source does not offer instead of stopping.")
	f.BoolVarP(&addYes, "yes", "y", false, "Apply without asking.")
	f.BoolVar(&addNoRefresh, "no-refresh", false, "Use Airbyte's cached source schema instead of rediscovering it.")
}
