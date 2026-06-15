package cmd

import (
	"abctlx/internal/airbyte"
	"fmt"

	"github.com/spf13/cobra"
)

var resCmd = &cobra.Command{
	Use:   "res",
	Short: "Interacts with Airbyte Resources",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("This command interacts with Airbyte Resources, --help for more info.")
	},
}

var updateConResRequest airbyte.UpdateConnectionResourceRequirementsRequest
var updateConResCmd = &cobra.Command{
	Use:   "update-connection",
	Short: "Edit Connection resources.",
	Run: func(cmd *cobra.Command, args []string) {
		//
	},
}

func init() {
	rootCmd.AddCommand(resCmd)
	resCmd.AddCommand(updateConResCmd)
	editConCmdFlags()
}

func editConCmdFlags() {
	updateConResCmd.Flags().StringVar(&updateConResRequest.ConnectionId, "id", "", "The Connection ID to update.")
	updateConResCmd.Flags().IntVar(&updateConResRequest.MinCpuCores, "min_cpu", 0, "The Minimum CPU usage")
	updateConResCmd.Flags().IntVar(&updateConResRequest.MaxCpuCores, "max_cpu", 0, "The Maximum CPU usage")
	updateConResCmd.Flags().IntVar(&updateConResRequest.MinMemGb, "min_mem", 0, "The Minimum RAM usage")
	updateConResCmd.Flags().IntVar(&updateConResRequest.MaxMemGb, "max_mem", 0, "The Maximum RAM usage")
}

// -c "UPDATE connection SET resource_requirements = '{\"cpu_limit\": \"2\", \"cpu_request\": \"2\", \"memory_limit\": \"2048Mi\", \"memory_request\": \"2048Mi\"}' WHERE id = 'YOUR_ID_HERE';"
