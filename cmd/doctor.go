package cmd

import (
	"abctlx/internal/airbyte"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the profile, secret, API, authentication and workspace",
	Long: `Walks from the profile file to a working API call and reports each step as
ok, warn, FAIL or skip, with a hint for anything that is not ok. Read-only: it
sends no change to Airbyte. Exits 1 if any check fails.

Uses the active profile, or --profile.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		checks := airbyte.Doctor(context.Background(), flagProfile)
		fails := 0
		fmt.Println()
		for _, c := range checks {
			fmt.Printf("  %-4s  %-15s %s\n", c.Status, c.Name, c.Detail)
			if c.Hint != "" && (c.Status == airbyte.CheckFail || c.Status == airbyte.CheckWarn) {
				fmt.Printf("        %-15s -> %s\n", "", c.Hint)
			}
			if c.Status == airbyte.CheckFail {
				fails++
			}
		}
		fmt.Println()
		if fails > 0 {
			cmd.SilenceUsage = true
			return fmt.Errorf("%d check(s) failed", fails)
		}
		fmt.Println("All checks passed.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
