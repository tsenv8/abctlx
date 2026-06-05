package cmd

import (
	"abctlx/internal/airbyte"
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var generateAccessTokenCmd = &cobra.Command{
	Use: "token",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		s := airbyte.NewAirbyteService(ctx)
		res, err := s.GenerateAccessToken()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf(res.AccessToken)
	},
}

func init() {
	rootCmd.AddCommand(generateAccessTokenCmd)
}
