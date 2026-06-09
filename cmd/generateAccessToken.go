package cmd

import (
	"abctlx/internal/airbyte"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var generateAccessTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generates an access token using Airbyte credentials.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		s := airbyte.NewAirbyteService(ctx)
		token := s.GetAccessToken()
		fmt.Printf("Token: %s", token)
	},
}

func init() {
	rootCmd.AddCommand(generateAccessTokenCmd)
}
