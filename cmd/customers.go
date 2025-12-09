package cmd

import (
	"time"

	"github.com/payjp/payjp-cli/internal/client"
	"github.com/payjp/payjp-cli/internal/util"
	"github.com/payjp/payjp-go/v1"
	"github.com/spf13/cobra"
)

var customersCmd = &cobra.Command{
	Use:     "customers",
	Aliases: []string{"cu", "customer"},
	Short:   "Manage customers",
	Long:    `Create, retrieve, update, and manage customers.`,
}

var customersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new customer",
	Long: `Create a new customer.

Example:
  payjp customers create --email user@example.com
  payjp customers create --email user@example.com --card tok_xxxxx
  payjp customers create --id my_customer_id --email user@example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")
		email, _ := cmd.Flags().GetString("email")
		description, _ := cmd.Flags().GetString("description")
		card, _ := cmd.Flags().GetString("card")
		metadata, _ := cmd.Flags().GetString("metadata")

		customer := payjp.Customer{}

		if id != "" {
			customer.ID = id
		}
		if email != "" {
			customer.Email = email
		}
		if description != "" {
			customer.Description = description
		}
		if card != "" {
			customer.CardToken = card
		}
		if metadata != "" {
			customer.Metadata = util.ParseMetadata(metadata)
		}

		result, err := client.GetCustomer().Create(customer)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var customersGetCmd = &cobra.Command{
	Use:   "get <customer_id>",
	Short: "Get customer information",
	Long: `Retrieve information about a specific customer.

Example:
  payjp customers get cus_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		customerID := args[0]

		result, err := client.GetCustomer().Retrieve(customerID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var customersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List customers",
	Long: `List all customers with optional filters.

Example:
  payjp customers list --limit 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		since, _ := cmd.Flags().GetString("since")
		until, _ := cmd.Flags().GetString("until")

		caller := client.GetCustomer().List()

		if limit > 0 {
			caller.Limit(limit)
		}
		if offset > 0 {
			caller.Offset(offset)
		}
		if since != "" {
			ts, err := util.ParseTimestamp(since)
			if err != nil {
				return err
			}
			caller.Since(time.Unix(ts, 0))
		}
		if until != "" {
			ts, err := util.ParseTimestamp(until)
			if err != nil {
				return err
			}
			caller.Until(time.Unix(ts, 0))
		}

		result, _, err := caller.Do()
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var customersUpdateCmd = &cobra.Command{
	Use:   "update <customer_id>",
	Short: "Update customer information",
	Long: `Update information for a specific customer.

Example:
  payjp customers update cus_xxxxx --email new@example.com
  payjp customers update cus_xxxxx --default-card car_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		customerID := args[0]
		email, _ := cmd.Flags().GetString("email")
		description, _ := cmd.Flags().GetString("description")
		defaultCard, _ := cmd.Flags().GetString("default-card")
		metadata, _ := cmd.Flags().GetString("metadata")

		customer := payjp.Customer{}

		if email != "" {
			customer.Email = email
		}
		if description != "" {
			customer.Description = description
		}
		if defaultCard != "" {
			customer.DefaultCard = defaultCard
		}
		if metadata != "" {
			customer.Metadata = util.ParseMetadata(metadata)
		}

		result, err := client.GetCustomer().Update(customerID, customer)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var customersDeleteCmd = &cobra.Command{
	Use:   "delete <customer_id>",
	Short: "Delete a customer",
	Long: `Delete a specific customer.

Example:
  payjp customers delete cus_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		customerID := args[0]

		err := client.GetCustomer().Delete(customerID)
		if err != nil {
			handleError(err)
			return nil
		}

		if quiet {
			return nil
		}

		return outputResult(map[string]interface{}{
			"id":      customerID,
			"deleted": true,
		})
	},
}

func init() {
	rootCmd.AddCommand(customersCmd)

	customersCmd.AddCommand(customersCreateCmd)
	customersCmd.AddCommand(customersGetCmd)
	customersCmd.AddCommand(customersListCmd)
	customersCmd.AddCommand(customersUpdateCmd)
	customersCmd.AddCommand(customersDeleteCmd)

	// Create flags
	customersCreateCmd.Flags().String("id", "", "Custom customer ID")
	customersCreateCmd.Flags().String("email", "", "Customer email")
	customersCreateCmd.Flags().String("description", "", "Description")
	customersCreateCmd.Flags().String("card", "", "Token ID to add as default card")
	customersCreateCmd.Flags().String("metadata", "", "Metadata (key1=value1,key2=value2)")

	// List flags
	customersListCmd.Flags().Int("limit", 10, "Number of items to return")
	customersListCmd.Flags().Int("offset", 0, "Offset for pagination")
	customersListCmd.Flags().String("since", "", "Filter by created timestamp (Unix timestamp or RFC3339)")
	customersListCmd.Flags().String("until", "", "Filter by created timestamp (Unix timestamp or RFC3339)")

	// Update flags
	customersUpdateCmd.Flags().String("email", "", "New email")
	customersUpdateCmd.Flags().String("description", "", "New description")
	customersUpdateCmd.Flags().String("default-card", "", "Card ID to set as default")
	customersUpdateCmd.Flags().String("metadata", "", "Metadata (key1=value1,key2=value2)")
}
