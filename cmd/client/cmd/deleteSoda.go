package cmd

import (
	v1 "colaco-api/internal/api/v1"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strings"
)

// deleteSodaCmd represents the deleteSoda command
var deleteSodaCmd = &cobra.Command{
	Use:   "delete-soda",
	Short: "deletes soda from the vending machine by removing the vending slot",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := v1.NewClientWithResponses(serverURL)
		if err != nil {
			log.Fatalf("couldn't create client with error: %v", err)
		}
		resp, err := client.AuthLoginWithResponse(context.Background(), v1.AuthLoginJSONRequestBody{
			Username: username,
			Password: password,
		})
		var token string
		if resp.JSON200 != nil {
			token = *resp.JSON200.Token
		} else {
			log.Println("Authentication failed or did not return a token")
		}

		authHeaderEditor := func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			return nil
		}

		soda, err := cmd.Flags().GetString("soda")
		if err != nil || soda == "" {
			log.Fatalf("soda name must be provided: %v", err)
		}
		soda = strings.ToLower(soda)

		r, err := client.DeleteVendingWithResponse(context.Background(), v1.DeleteVendingJSONRequestBody{
			Name: strings.ToLower(soda),
		}, authHeaderEditor)
		if r.JSON200 != nil {
			fmt.Println("soda deleted successfully")
		} else if r.JSON404 != nil {
			fmt.Printf("soda not found: %v\n", soda)
		} else {
			fmt.Println("something went wrong")
		}

	},
}

func init() {
	rootCmd.AddCommand(deleteSodaCmd)
	deleteSodaCmd.Flags().StringP("soda", "", "", "soda to delete from the vending machine")
}
