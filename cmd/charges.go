package cmd

import (
	"time"

	"github.com/payjp/payjp-cli/internal/client"
	"github.com/payjp/payjp-cli/internal/util"
	"github.com/payjp/payjp-go/v1"
	"github.com/spf13/cobra"
)

var chargesCmd = &cobra.Command{
	Use:     "charges",
	Aliases: []string{"ch", "charge"},
	Short:   "Manage charges",
	Long:    `Create, retrieve, update, and manage payment charges.`,
}

var chargesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new charge",
	Long: `Create a new charge (payment).

Example:
  payjp charges create --amount 1000 --currency jpy --card tok_xxxxx
  payjp charges create --amount 1000 --currency jpy --customer cus_xxxxx
  payjp charges create --amount 1000 --currency jpy --card tok_xxxxx --capture=false`,
	RunE: func(cmd *cobra.Command, args []string) error {
		amount, _ := cmd.Flags().GetInt("amount")
		currency, _ := cmd.Flags().GetString("currency")
		card, _ := cmd.Flags().GetString("card")
		customer, _ := cmd.Flags().GetString("customer")
		description, _ := cmd.Flags().GetString("description")
		capture, _ := cmd.Flags().GetBool("capture")
		expiryDays, _ := cmd.Flags().GetInt("expiry-days")
		metadata, _ := cmd.Flags().GetString("metadata")
		threeDSecure, _ := cmd.Flags().GetBool("three-d-secure")

		if err := util.ValidateAmount(amount); err != nil {
			return err
		}
		if err := util.ValidateCurrency(currency); err != nil {
			return err
		}

		charge := payjp.Charge{
			Currency: currency,
			Capture:  capture,
		}

		if card != "" {
			charge.CardToken = card
		}
		if customer != "" {
			charge.CustomerID = customer
		}
		if description != "" {
			charge.Description = description
		}
		if expiryDays > 0 {
			charge.ExpireDays = expiryDays
		}
		if metadata != "" {
			charge.Metadata = util.ParseMetadata(metadata)
		}
		if threeDSecure {
			tds := true
			charge.ThreeDSecure = &tds
		}

		result, err := client.GetCharge().Create(amount, charge)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var chargesGetCmd = &cobra.Command{
	Use:   "get <charge_id>",
	Short: "Get charge information",
	Long: `Retrieve information about a specific charge.

Example:
  payjp charges get ch_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		chargeID := args[0]

		result, err := client.GetCharge().Retrieve(chargeID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var chargesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List charges",
	Long: `List all charges with optional filters.

Example:
  payjp charges list --limit 10
  payjp charges list --customer cus_xxxxx`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		since, _ := cmd.Flags().GetString("since")
		until, _ := cmd.Flags().GetString("until")
		customer, _ := cmd.Flags().GetString("customer")
		subscription, _ := cmd.Flags().GetString("subscription")

		caller := client.GetCharge().List()

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
		if customer != "" {
			caller.CustomerID(customer)
		}
		if subscription != "" {
			caller.SubscriptionID(subscription)
		}

		result, _, err := caller.Do()
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var chargesUpdateCmd = &cobra.Command{
	Use:   "update <charge_id>",
	Short: "Update charge information",
	Long: `Update information for a specific charge.

Example:
  payjp charges update ch_xxxxx --description "New description"
  payjp charges update ch_xxxxx --metadata key1=value1`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		chargeID := args[0]
		description, _ := cmd.Flags().GetString("description")
		metadata, _ := cmd.Flags().GetString("metadata")

		var result *payjp.ChargeResponse
		var err error

		if metadata != "" {
			result, err = client.GetCharge().Update(chargeID, description, util.ParseMetadata(metadata))
		} else {
			result, err = client.GetCharge().Update(chargeID, description)
		}

		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var chargesCaptureCmd = &cobra.Command{
	Use:   "capture <charge_id>",
	Short: "Capture an authorized charge",
	Long: `Capture an authorized charge.

Example:
  payjp charges capture ch_xxxxx
  payjp charges capture ch_xxxxx --amount 500`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		chargeID := args[0]
		amount, _ := cmd.Flags().GetInt("amount")

		var result *payjp.ChargeResponse
		var err error

		if amount > 0 {
			result, err = client.GetCharge().Capture(chargeID, amount)
		} else {
			result, err = client.GetCharge().Capture(chargeID)
		}

		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var chargesRefundCmd = &cobra.Command{
	Use:   "refund <charge_id>",
	Short: "Refund a charge",
	Long: `Refund a captured charge.

Example:
  payjp charges refund ch_xxxxx
  payjp charges refund ch_xxxxx --amount 500
  payjp charges refund ch_xxxxx --refund-reason "Customer request"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		chargeID := args[0]
		amount, _ := cmd.Flags().GetInt("amount")
		refundReason, _ := cmd.Flags().GetString("refund-reason")

		var result *payjp.ChargeResponse
		var err error

		// SDK signature: Refund(chargeID, reason string, amount ...int)
		if amount > 0 {
			result, err = client.GetCharge().Refund(chargeID, refundReason, amount)
		} else {
			result, err = client.GetCharge().Refund(chargeID, refundReason)
		}

		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var chargesTdsFinishCmd = &cobra.Command{
	Use:   "tds-finish <charge_id>",
	Short: "Complete 3D Secure authentication",
	Long: `Complete 3D Secure authentication for a charge.

Example:
  payjp charges tds-finish ch_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		chargeID := args[0]

		result, err := client.GetCharge().TdsFinish(chargeID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

func init() {
	rootCmd.AddCommand(chargesCmd)

	chargesCmd.AddCommand(chargesCreateCmd)
	chargesCmd.AddCommand(chargesGetCmd)
	chargesCmd.AddCommand(chargesListCmd)
	chargesCmd.AddCommand(chargesUpdateCmd)
	chargesCmd.AddCommand(chargesCaptureCmd)
	chargesCmd.AddCommand(chargesRefundCmd)
	chargesCmd.AddCommand(chargesTdsFinishCmd)

	// Create flags
	chargesCreateCmd.Flags().Int("amount", 0, "Amount in smallest currency unit (required)")
	chargesCreateCmd.Flags().String("currency", "jpy", "Currency code")
	chargesCreateCmd.Flags().String("card", "", "Token ID")
	chargesCreateCmd.Flags().String("customer", "", "Customer ID")
	chargesCreateCmd.Flags().String("description", "", "Description")
	chargesCreateCmd.Flags().Bool("capture", true, "Capture immediately")
	chargesCreateCmd.Flags().Int("expiry-days", 0, "Expiry days for authorization")
	chargesCreateCmd.Flags().String("metadata", "", "Metadata (key1=value1,key2=value2)")
	chargesCreateCmd.Flags().Bool("three-d-secure", false, "Enable 3D Secure")
	chargesCreateCmd.MarkFlagRequired("amount")

	// List flags
	chargesListCmd.Flags().Int("limit", 10, "Number of items to return")
	chargesListCmd.Flags().Int("offset", 0, "Offset for pagination")
	chargesListCmd.Flags().String("since", "", "Filter by created timestamp (Unix timestamp or RFC3339)")
	chargesListCmd.Flags().String("until", "", "Filter by created timestamp (Unix timestamp or RFC3339)")
	chargesListCmd.Flags().String("customer", "", "Filter by customer ID")
	chargesListCmd.Flags().String("subscription", "", "Filter by subscription ID")

	// Update flags
	chargesUpdateCmd.Flags().String("description", "", "New description")
	chargesUpdateCmd.Flags().String("metadata", "", "Metadata (key1=value1,key2=value2)")

	// Capture flags
	chargesCaptureCmd.Flags().Int("amount", 0, "Amount to capture (partial capture)")

	// Refund flags
	chargesRefundCmd.Flags().Int("amount", 0, "Amount to refund (partial refund)")
	chargesRefundCmd.Flags().String("refund-reason", "", "Reason for refund")
}

