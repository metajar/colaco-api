package cmd

import (
	v1 "colaco-api/internal/api/v1"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

var addSodaCmd = &cobra.Command{
	Use:   "add-soda",
	Short: "Adds a new soda to the vending machine",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := v1.NewClientWithResponses(serverURL)
		if err != nil {
			log.Fatalf("couldn't create client with error: %v", err)
		}
		token, err := authenticate(client)
		if err != nil {
			log.Fatalf("authentication failed: %v", err)
		}

		sodaName, err := cmd.Flags().GetString("name")
		if err != nil || sodaName == "" {
			log.Fatalf("name must be provided: %v", err)
		}
		price, err := cmd.Flags().GetFloat32("price")
		if err != nil {
			log.Fatalf("price must be provided: %v", err)
		}
		quantity, err := cmd.Flags().GetInt("quantity")
		if err != nil {
			log.Fatalf("quantity must be provided: %v", err)
		}
		calories, err := cmd.Flags().GetInt("calories")
		if err != nil {
			log.Fatalf("calories must be provided: %v", err)
		}
		ounces, err := cmd.Flags().GetFloat32("ounces")
		if err != nil {
			log.Fatalf("ounces must be provided: %v", err)
		}
		description, err := cmd.Flags().GetString("description")
		if err != nil {
			log.Fatalf("description must be provided: %v", err)
		}
		origin, err := cmd.Flags().GetString("origin")
		if err != nil {
			log.Fatalf("origin must be provided: %v", err)
		}

		newSoda := v1.PostNewJSONRequestBody{
			Slot: v1.VendingSlot{
				Cost:     &price,
				Quantity: &quantity,
				OccupiedSoda: &v1.Soda{
					Name:        &sodaName,
					Calories:    &calories,
					Ounces:      &ounces,
					Description: &description,
					OriginStory: &origin,
				},
			},
		}

		r, err := client.PostNewWithResponse(context.Background(), newSoda, func(ctx context.Context, req *http.Request) error {
			return addAuthHeader(ctx, req, token)
		})
		if err != nil {
			log.Fatalf("Failed to add new soda: %v", err)
		}
		if r.JSON201 != nil {
			fmt.Println("Soda added successfully")
		} else if r.JSON409 != nil {
			fmt.Println("Soda conflict found. Already exists.")
		} else {
			fmt.Println("An unexpected error occurred")
		}
	},
}

func init() {
	rootCmd.AddCommand(addSodaCmd)
	addSodaCmd.Flags().StringP("name", "", "", "Name of the soda")
	addSodaCmd.Flags().StringP("origin", "", "", "Origin story of the soda")
	addSodaCmd.Flags().StringP("description", "", "", "Description of the soda")
	addSodaCmd.Flags().Float32P("price", "", 0.0, "Price of the soda")
	addSodaCmd.Flags().IntP("quantity", "", 0, "Initial quantity of the soda")
	addSodaCmd.Flags().IntP("calories", "", 0, "Calories of the soda")
	addSodaCmd.Flags().Float32P("ounces", "", 0.0, "Ounces of the soda")
	addSodaCmd.MarkFlagRequired("name")
	addSodaCmd.MarkFlagRequired("price")
	addSodaCmd.MarkFlagRequired("quantity")
}
