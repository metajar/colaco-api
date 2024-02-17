package cmd

import (
	v1 "colaco-api/internal/api/v1"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// authLoginCmd represents the authLogin command
var authLoginCmd = &cobra.Command{
	Use:   "get-token",
	Short: "gets token from the server that can be used with other tooling such as postman.",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := v1.NewClientWithResponses(serverURL)
		if err != nil {
			log.Fatalf("couldn't create client with error: %v", err)
		}
		resp, err := client.AuthLoginWithResponse(context.Background(), v1.AuthLoginJSONRequestBody{
			Username: username,
			Password: password,
		})
		if resp.JSON200 != nil {
			token := resp.JSON200.Token
			fmt.Printf("Token: %s\n", *token)
		} else {
			log.Println("Authentication failed or did not return a token")
		}

	},
}

func init() {
	rootCmd.AddCommand(authLoginCmd)
}
