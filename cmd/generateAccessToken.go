package cmd

import (
	"abctlx/internal/airbyte"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var generateAccessTokenCmd = &cobra.Command{
	Use: "token",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		s := airbyte.NewAirbyteService(ctx)
		res := s.GenerateAccessToken()
		fmt.Printf("Token: %s", res.AccessToken)
	},
}

func init() {
	rootCmd.AddCommand(generateAccessTokenCmd)
}
