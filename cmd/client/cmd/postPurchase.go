package cmd

import (
	v1 "colaco-api/internal/api/v1"
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
)

var purchaseSodaCmd = &cobra.Command{
	Use:   "purchase-soda",
	Short: "Purchases a soda from the vending machine",
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
		payment, err := cmd.Flags().GetFloat32("payment")
		if err != nil {
			log.Fatalf("payment must be provided: %v", err)
		}

		purchaseRequest := v1.PostPurchaseJSONRequestBody{
			Name:    sodaName,
			Payment: payment,
		}

		r, err := client.PostPurchaseWithResponse(context.Background(), purchaseRequest, func(ctx context.Context, req *http.Request) error {
			return addAuthHeader(ctx, req, token)
		})
		if err != nil {
			log.Fatalf("Failed to purchase soda: %v", err)
		}
		if r.JSON200 != nil {
			displayPurchaseDetails(r.JSON200)
			fmt.Println("\nEnjoy your drink!")
		} else if r.JSON402 != nil {
			fmt.Printf("Insufficient funds. Please add more funds.")
		} else {
			fmt.Println("An unexpected error occurred")
		}
	},
}

func init() {
	rootCmd.AddCommand(purchaseSodaCmd)
	purchaseSodaCmd.Flags().StringP("soda", "", "", "Name of the soda to purchase")
	purchaseSodaCmd.Flags().Float32P("payment", "", 0.0, "Payment amount")
	purchaseSodaCmd.MarkFlagRequired("soda")
	purchaseSodaCmd.MarkFlagRequired("payment")
}

func displayPurchaseDetails(details *v1.PurchaseSodaResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Attribute", "Details"})
	table.SetBorder(true)
	table.SetColumnSeparator(":")
	table.Append([]string{"Soda Name", *details.Soda.Name})
	table.Append([]string{"Description", *details.Soda.Description})
	table.Append([]string{"Calories", fmt.Sprintf("%d", *details.Soda.Calories)})
	table.Append([]string{"Volume (Ounces)", fmt.Sprintf("%.1f", *details.Soda.Ounces)})
	table.Append([]string{"Origin Story", *details.Soda.OriginStory})
	table.Append([]string{"Change Returned", fmt.Sprintf("$%.2f", *details.Change)})

	fmt.Println("Dispensing your soda...")
	table.Render() // Print the table to the console
}
