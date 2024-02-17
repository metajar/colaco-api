package cmd

import (
	v1 "colaco-api/internal/api/v1"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var getVendingCmd = &cobra.Command{
	Use:   "get-sodas",
	Short: "Gathers all the sodas that are in the vending slots.",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := v1.NewClientWithResponses(serverURL)
		if err != nil {
			log.Fatalf("couldn't create client with error: %v", err)
		}
		token, err := authenticate(client)
		if err != nil {
			log.Fatalf("authentication failed: %v", err)
		}

		r, err := client.GetVendingWithResponse(context.Background(), v1.GetVendingJSONRequestBody{Name: ""}, func(ctx context.Context, req *http.Request) error {
			return addAuthHeader(ctx, req, token)
		})
		if err != nil {
			log.Fatalf("Failed to get sodas: %v", err)
		}

		if r.JSON200 != nil {
			printSodaTable(*r.JSON200.Slots)
		} else if r.JSON404 != nil {
			fmt.Println("No sodas found")
		} else {
			fmt.Println("An unexpected error occurred")
		}
	},
}

func init() {
	rootCmd.AddCommand(getVendingCmd)
}

func printSodaTable(sodas []v1.VendingSlot) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "Soda Name\tCalories\tOunces\tCost\tQuantity\tDescription")
	for _, soda := range sodas {
		if soda.OccupiedSoda != nil {
			fmt.Fprintf(w, "%s\t%d\t%.2f\t%.2f\t%d\t%s\n",
				*soda.OccupiedSoda.Name,
				*soda.OccupiedSoda.Calories,
				*soda.OccupiedSoda.Ounces,
				*soda.Cost,
				*soda.Quantity,
				*soda.OccupiedSoda.Description,
			)
		}
	}
	w.Flush()
}
