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

var updatePriceCmd = &cobra.Command{
	Use:   "update-price",
	Short: "updates the price of a soda",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := v1.NewClientWithResponses(serverURL)
		if err != nil {
			log.Fatalf("couldn't create client with error: %v", err)
		}
		token, err := authenticate(client)
		if err != nil {
			log.Fatalf("authentication failed: %v", err)
		}

		soda, err := cmd.Flags().GetString("soda")
		if err != nil || soda == "" {
			log.Fatalf("soda name must be provided: %v", err)
		}
		soda = strings.ToLower(soda)

		price, err := cmd.Flags().GetFloat32("price")
		if err != nil || price <= 0 { // Assuming price must be greater than 0
			log.Fatalf("valid soda price must be provided: %v", err)
		}

		r, err := client.UpdatePriceWithResponse(context.Background(), v1.UpdatePriceJSONRequestBody{
			Name:     soda,
			NewPrice: price,
		}, func(ctx context.Context, req *http.Request) error {
			return addAuthHeader(ctx, req, token)
		})
		if err != nil {
			log.Fatalf("Failed to update soda price: %v", err)
		}
		if r.JSON200 != nil {
			fmt.Printf("Soda price updated successfully from %v to %v.\n", *r.JSON200.OldPrice, *r.JSON200.NewPrice)
		} else if r.JSON404 != nil {
			fmt.Printf("Soda not found: %v\n", soda)
		} else {
			fmt.Println("An unexpected error occurred")
		}
	},
}

func init() {
	rootCmd.AddCommand(updatePriceCmd)
	updatePriceCmd.Flags().StringP("soda", "", "", "Soda to update price on")
	updatePriceCmd.Flags().Float32P("price", "", 1.00, "Price to update soda to")
	updatePriceCmd.MarkFlagRequired("soda")
	updatePriceCmd.MarkFlagRequired("price")
}
