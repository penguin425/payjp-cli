package cmd

import (
	"fmt"

	"github.com/payjp/payjp-cli/internal/client"
	"github.com/payjp/payjp-cli/internal/util"
	"github.com/payjp/payjp-go/v1"
	"github.com/spf13/cobra"
)

var cardsCmd = &cobra.Command{
	Use:     "cards",
	Aliases: []string{"card"},
	Short:   "Manage customer cards",
	Long:    `Add, retrieve, update, and delete cards for customers.`,
}

var cardsCreateCmd = &cobra.Command{
	Use:   "create <customer_id>",
	Short: "Add a card to a customer",
	Long: `Add a new card to a customer using a token.

Example:
  payjp cards create cus_xxxxx --card tok_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		customerID := args[0]
		card, _ := cmd.Flags().GetString("card")

		if card == "" {
			return fmt.Errorf("--card is required")
		}

		customer, err := client.GetCustomer().Retrieve(customerID)
		if err != nil {
			handleError(err)
			return nil
		}

		result, err := customer.AddCardToken(card)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var cardsGetCmd = &cobra.Command{
	Use:   "get <customer_id> <card_id>",
	Short: "Get card information",
	Long: `Retrieve information about a specific card.

Example:
  payjp cards get cus_xxxxx car_xxxxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		customerID := args[0]
		cardID := args[1]

		customer, err := client.GetCustomer().Retrieve(customerID)
		if err != nil {
			handleError(err)
			return nil
		}

		result, err := customer.GetCard(cardID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var cardsListCmd = &cobra.Command{
	Use:   "list <customer_id>",
	Short: "List customer cards",
	Long: `List all cards for a specific customer.

Example:
  payjp cards list cus_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		customerID := args[0]
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		customer, err := client.GetCustomer().Retrieve(customerID)
		if err != nil {
			handleError(err)
			return nil
		}

		caller := customer.ListCard()
		if limit > 0 {
			caller.Limit(limit)
		}
		if offset > 0 {
			caller.Offset(offset)
		}

		result, _, err := caller.Do()
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var cardsUpdateCmd = &cobra.Command{
	Use:   "update <customer_id> <card_id>",
	Short: "Update card information",
	Long: `Update information for a specific card.

Example:
  payjp cards update cus_xxxxx car_xxxxx --name "PAY TARO"
  payjp cards update cus_xxxxx car_xxxxx --address-zip "1000001"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		customerID := args[0]
		cardID := args[1]
		name, _ := cmd.Flags().GetString("name")
		addressZip, _ := cmd.Flags().GetString("address-zip")
		addressState, _ := cmd.Flags().GetString("address-state")
		addressCity, _ := cmd.Flags().GetString("address-city")
		addressLine1, _ := cmd.Flags().GetString("address-line1")
		addressLine2, _ := cmd.Flags().GetString("address-line2")
		country, _ := cmd.Flags().GetString("country")
		metadata, _ := cmd.Flags().GetString("metadata")

		customer, err := client.GetCustomer().Retrieve(customerID)
		if err != nil {
			handleError(err)
			return nil
		}

		card := payjp.Card{}

		if name != "" {
			card.Name = name
		}
		if addressZip != "" {
			card.AddressZip = addressZip
		}
		if addressState != "" {
			card.AddressState = addressState
		}
		if addressCity != "" {
			card.AddressCity = addressCity
		}
		if addressLine1 != "" {
			card.AddressLine1 = addressLine1
		}
		if addressLine2 != "" {
			card.AddressLine2 = addressLine2
		}
		if country != "" {
			card.Country = country
		}
		if metadata != "" {
			card.Metadata = util.ParseMetadata(metadata)
		}

		result, err := customer.UpdateCard(cardID, card)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var cardsDeleteCmd = &cobra.Command{
	Use:   "delete <customer_id> <card_id>",
	Short: "Delete a card",
	Long: `Delete a specific card from a customer.

Example:
  payjp cards delete cus_xxxxx car_xxxxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		customerID := args[0]
		cardID := args[1]

		customer, err := client.GetCustomer().Retrieve(customerID)
		if err != nil {
			handleError(err)
			return nil
		}

		err = customer.DeleteCard(cardID)
		if err != nil {
			handleError(err)
			return nil
		}

		if quiet {
			return nil
		}

		return outputResult(map[string]interface{}{
			"id":      cardID,
			"deleted": true,
		})
	},
}

func init() {
	rootCmd.AddCommand(cardsCmd)

	cardsCmd.AddCommand(cardsCreateCmd)
	cardsCmd.AddCommand(cardsGetCmd)
	cardsCmd.AddCommand(cardsListCmd)
	cardsCmd.AddCommand(cardsUpdateCmd)
	cardsCmd.AddCommand(cardsDeleteCmd)

	// Create flags
	cardsCreateCmd.Flags().String("card", "", "Token ID (required)")
	cardsCreateCmd.MarkFlagRequired("card")

	// List flags
	cardsListCmd.Flags().Int("limit", 10, "Number of items to return")
	cardsListCmd.Flags().Int("offset", 0, "Offset for pagination")

	// Update flags
	cardsUpdateCmd.Flags().String("name", "", "Cardholder name")
	cardsUpdateCmd.Flags().String("address-zip", "", "Postal code")
	cardsUpdateCmd.Flags().String("address-state", "", "State/Prefecture")
	cardsUpdateCmd.Flags().String("address-city", "", "City")
	cardsUpdateCmd.Flags().String("address-line1", "", "Address line 1")
	cardsUpdateCmd.Flags().String("address-line2", "", "Address line 2")
	cardsUpdateCmd.Flags().String("country", "", "Country code (e.g., JP)")
	cardsUpdateCmd.Flags().String("metadata", "", "Metadata (key1=value1,key2=value2)")
}
