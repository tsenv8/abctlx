package cmd

import (
	"abctlx/helpers"
	"abctlx/internal/airbyte"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	addTag         string
	addStreamsId   string
	addSkipMissing bool
	addYes         bool
	addNoRefresh   bool
)

var addStreamsCmd = &cobra.Command{
	Use:   "add-streams (--id <conn> | --tag <tag>) <schema.table>...",
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
		if (addStreamsId == "") == (addTag == "") {
			return fmt.Errorf("pass exactly one of --id or --tag")
		}
		s := airbyte.NewAirbyteService(context.Background())

		var ids []string
		if addStreamsId != "" {
			ids = []string{addStreamsId}
		} else {
			limit := 500
			for _, c := range s.ListConnections(&limit).Data {
				for _, t := range c.Tags {
					if t.Name == addTag {
						ids = append(ids, c.ConnectionId)
						break
					}
				}
			}
			if len(ids) == 0 {
				return fmt.Errorf("no connection has the tag %q", addTag)
			}
		}

		// Plan every connection first; nothing is sent until all are planned.
		type target struct {
			con     airbyte.ConnectionData
			plan    airbyte.PlanResult
			streams []json.RawMessage
		}
		var targets []target
		var blocked []string
		for _, id := range ids {
			con, _, err := s.GetConnectionByID(id)
			if err != nil {
				return fmt.Errorf("get connection %s: %w", id, err)
			}
			avail, _, err := s.ListAvailableStreams(con.SourceId, con.DestinationId, !addNoRefresh)
			if err != nil {
				return fmt.Errorf("list available streams for %s: %w", con.Name, err)
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
				fmt.Sprintf("%s: %d enabled -> %d after", con.Name, len(con.Configurations.Streams), len(streams)))
			if len(plan.Already) > 0 {
				fmt.Printf("\nAlready enabled (skipped): %s\n", strings.Join(plan.Already, ", "))
			}
			if len(plan.Missing) > 0 {
				fmt.Printf("\nNot offered by the source: %s\n", strings.Join(plan.Missing, ", "))
				if !addSkipMissing {
					blocked = append(blocked, fmt.Sprintf("%s (%s)", con.Name, strings.Join(plan.Missing, ", ")))
				}
			}
			targets = append(targets, target{con, plan, streams})
		}
		if len(blocked) > 0 {
			return fmt.Errorf("nothing changed: tables missing on %s (publish them first, or pass --skip-missing)", strings.Join(blocked, "; "))
		}

		var work []target
		for _, t := range targets {
			if len(t.plan.Add) > 0 {
				work = append(work, t)
			}
		}
		if len(work) == 0 {
			fmt.Println("\nNothing to add.")
			return nil
		}
		if airbyte.DryRun {
			fmt.Printf("\nDry run: nothing sent (%d connection(s) would change).\n", len(work))
			return nil
		}
		if !addYes {
			var names []string
			for _, t := range work {
				names = append(names, t.con.Name)
			}
			fmt.Printf("\nApply to %d connection(s): %s? [y/N] ", len(work), strings.Join(names, ", "))
			ans, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			if a := strings.ToLower(strings.TrimSpace(ans)); a != "y" && a != "yes" {
				fmt.Println("Cancelled; nothing changed.")
				return nil
			}
		}
		var failed []string
		for _, t := range work {
			updated, err := s.UpdateConnectionStreams(t.con.ConnectionId, t.streams)
			switch {
			case err != nil:
				failed = append(failed, fmt.Sprintf("%s: %v", t.con.Name, err))
			case len(updated.Configurations.Streams) != len(t.streams):
				failed = append(failed, fmt.Sprintf("%s: expected %d streams, Airbyte reports %d", t.con.Name, len(t.streams), len(updated.Configurations.Streams)))
			default:
				fmt.Printf("\nDone: %s now has %d streams.\n", t.con.Name, len(updated.Configurations.Streams))
			}
		}
		if len(failed) > 0 {
			return fmt.Errorf("failed: %s", strings.Join(failed, "; "))
		}
		return nil
	},
}

func init() {
	conCmd.AddCommand(addStreamsCmd)
	f := addStreamsCmd.Flags()
	f.StringVar(&addStreamsId, "id", "", "Connection ID.")
	f.StringVar(&addTag, "tag", "", "Every connection with this tag (e.g. client), planned together; nothing is sent unless all are clean.")
	f.BoolVar(&addSkipMissing, "skip-missing", false, "Skip tables the source does not offer instead of stopping.")
	f.BoolVarP(&addYes, "yes", "y", false, "Apply without asking.")
	f.BoolVar(&addNoRefresh, "no-refresh", false, "Use Airbyte's cached source schema instead of rediscovering it.")
}
