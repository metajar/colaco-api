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

var restockSodaCmd = &cobra.Command{
	Use:   "restock-soda",
	Short: "Restocks a specific soda in the vending machine",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := v1.NewClientWithResponses(serverURL)
		if err != nil {
			log.Fatalf("couldn't create client with error: %v", err)
		}
		token, err := authenticate(client)
		if err != nil {
			log.Fatalf("authentication failed: %v", err)
		}

		sodaName, err := cmd.Flags().GetString("soda")
		if err != nil || sodaName == "" {
			log.Fatalf("soda name must be provided: %v", err)
		}
		sodaName = strings.ToLower(sodaName)

		quantity, err := cmd.Flags().GetInt("qty")
		if err != nil || quantity <= 0 {
			log.Fatalf("quantity must be provided and greater than 0: %v", err)
		}

		r, err := client.RestockSodaWithResponse(context.Background(), v1.RestockSodaJSONRequestBody{
			Name:     sodaName,
			Quantity: quantity,
		}, func(ctx context.Context, req *http.Request) error {
			return addAuthHeader(ctx, req, token)
		})
		if err != nil {
			log.Fatalf("failed to restock soda: %v", err)
		}

		if r.JSON200 != nil {
			fmt.Printf("Successfully restocked %s with %d additional units.\n", sodaName, quantity)
			if r.JSON200.Leftover != nil && *r.JSON200.Leftover > 0 {
				fmt.Printf("Warning: %d units could not be added due to capacity limits.\n", *r.JSON200.Leftover)
			}
		} else if r.JSON404 != nil {
			fmt.Printf("Soda '%s' not found.\n", sodaName)
		} else {
			fmt.Println("An unexpected error occurred.")
		}
	},
}

func init() {
	rootCmd.AddCommand(restockSodaCmd)
	restockSodaCmd.Flags().StringP("soda", "", "", "Soda to replenish in the vending machine")
	restockSodaCmd.Flags().IntP("qty", "", 0, "Quantity of soda to add")
	// Ensuring the necessary flags are marked as required
	restockSodaCmd.MarkFlagRequired("soda")
	restockSodaCmd.MarkFlagRequired("qty")
}
