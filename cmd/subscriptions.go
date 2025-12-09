package cmd

import (
	"fmt"
	"time"

	"github.com/payjp/payjp-cli/internal/client"
	"github.com/payjp/payjp-cli/internal/util"
	"github.com/payjp/payjp-go/v1"
	"github.com/spf13/cobra"
)

var subscriptionsCmd = &cobra.Command{
	Use:     "subscriptions",
	Aliases: []string{"sub", "subscription"},
	Short:   "Manage subscriptions",
	Long:    `Create, retrieve, update, and manage recurring subscriptions.`,
}

var subscriptionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new subscription",
	Long: `Create a new recurring subscription.

Example:
  payjp subscriptions create --customer cus_xxxxx --plan pln_xxxxx
  payjp subscriptions create --customer cus_xxxxx --plan pln_xxxxx --trial-end 1640000000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		customer, _ := cmd.Flags().GetString("customer")
		plan, _ := cmd.Flags().GetString("plan")
		trialEnd, _ := cmd.Flags().GetString("trial-end")
		prorate, _ := cmd.Flags().GetBool("prorate")
		metadata, _ := cmd.Flags().GetString("metadata")

		subscription := payjp.Subscription{
			PlanID:  plan,
			Prorate: prorate,
		}

		if trialEnd != "" {
			ts, err := util.ParseTimestamp(trialEnd)
			if err != nil {
				return err
			}
			subscription.TrialEnd = time.Unix(ts, 0)
		}
		if metadata != "" {
			subscription.Metadata = util.ParseMetadata(metadata)
		}

		result, err := client.GetSubscription().Subscribe(customer, subscription)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var subscriptionsGetCmd = &cobra.Command{
	Use:   "get <customer_id> <subscription_id>",
	Short: "Get subscription information",
	Long: `Retrieve information about a specific subscription.

Example:
  payjp subscriptions get cus_xxxxx sub_xxxxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		customerID := args[0]
		subscriptionID := args[1]

		result, err := client.GetSubscription().Retrieve(customerID, subscriptionID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var subscriptionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List subscriptions",
	Long: `List all subscriptions with optional filters.

Example:
  payjp subscriptions list --limit 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		caller := client.GetSubscription().List()

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

var subscriptionsUpdateCmd = &cobra.Command{
	Use:   "update <subscription_id>",
	Short: "Update subscription information",
	Long: `Update information for a specific subscription.

Example:
  payjp subscriptions update sub_xxxxx --plan pln_new_xxxxx
  payjp subscriptions update sub_xxxxx --trial-end 1640000000`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]
		plan, _ := cmd.Flags().GetString("plan")
		trialEnd, _ := cmd.Flags().GetString("trial-end")
		prorate, _ := cmd.Flags().GetBool("prorate")
		metadata, _ := cmd.Flags().GetString("metadata")

		subscription := payjp.Subscription{
			Prorate: prorate,
		}

		if plan != "" {
			subscription.PlanID = plan
		}
		if trialEnd != "" {
			ts, err := util.ParseTimestamp(trialEnd)
			if err != nil {
				return err
			}
			subscription.TrialEnd = time.Unix(ts, 0)
		}
		if metadata != "" {
			subscription.Metadata = util.ParseMetadata(metadata)
		}

		result, err := client.GetSubscription().Update(subscriptionID, subscription)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var subscriptionsPauseCmd = &cobra.Command{
	Use:   "pause <subscription_id>",
	Short: "Pause a subscription",
	Long: `Pause an active subscription.

Example:
  payjp subscriptions pause sub_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		result, err := client.GetSubscription().Pause(subscriptionID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var subscriptionsResumeCmd = &cobra.Command{
	Use:   "resume <subscription_id>",
	Short: "Resume a paused subscription",
	Long: `Resume a paused subscription.

Example:
  payjp subscriptions resume sub_xxxxx
  payjp subscriptions resume sub_xxxxx --trial-end 1640000000`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]
		trialEnd, _ := cmd.Flags().GetString("trial-end")
		prorate, _ := cmd.Flags().GetBool("prorate")

		subscription := payjp.Subscription{
			Prorate: prorate,
		}

		if trialEnd != "" {
			ts, err := util.ParseTimestamp(trialEnd)
			if err != nil {
				return err
			}
			subscription.TrialEnd = time.Unix(ts, 0)
		}

		result, err := client.GetSubscription().Resume(subscriptionID, subscription)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var subscriptionsCancelCmd = &cobra.Command{
	Use:   "cancel <subscription_id>",
	Short: "Cancel a subscription",
	Long: `Cancel a subscription at the end of the current period.

Example:
  payjp subscriptions cancel sub_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		result, err := client.GetSubscription().Cancel(subscriptionID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var subscriptionsDeleteCmd = &cobra.Command{
	Use:   "delete <subscription_id>",
	Short: "Delete a subscription",
	Long: `Delete a subscription immediately.

Example:
  payjp subscriptions delete sub_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		err := client.GetSubscription().Delete(subscriptionID, payjp.SubscriptionDelete{})
		if err != nil {
			handleError(err)
			return nil
		}

		if quiet {
			fmt.Println(subscriptionID)
			return nil
		}

		return outputResult(map[string]interface{}{
			"id":      subscriptionID,
			"deleted": true,
		})
	},
}

func init() {
	rootCmd.AddCommand(subscriptionsCmd)

	subscriptionsCmd.AddCommand(subscriptionsCreateCmd)
	subscriptionsCmd.AddCommand(subscriptionsGetCmd)
	subscriptionsCmd.AddCommand(subscriptionsListCmd)
	subscriptionsCmd.AddCommand(subscriptionsUpdateCmd)
	subscriptionsCmd.AddCommand(subscriptionsPauseCmd)
	subscriptionsCmd.AddCommand(subscriptionsResumeCmd)
	subscriptionsCmd.AddCommand(subscriptionsCancelCmd)
	subscriptionsCmd.AddCommand(subscriptionsDeleteCmd)

	// Create flags
	subscriptionsCreateCmd.Flags().String("customer", "", "Customer ID (required)")
	subscriptionsCreateCmd.Flags().String("plan", "", "Plan ID (required)")
	subscriptionsCreateCmd.Flags().String("trial-end", "", "Trial end timestamp (Unix timestamp or RFC3339)")
	subscriptionsCreateCmd.Flags().Bool("prorate", false, "Prorate charges")
	subscriptionsCreateCmd.Flags().String("metadata", "", "Metadata (key1=value1,key2=value2)")
	subscriptionsCreateCmd.MarkFlagRequired("customer")
	subscriptionsCreateCmd.MarkFlagRequired("plan")

	// List flags
	subscriptionsListCmd.Flags().Int("limit", 10, "Number of items to return")
	subscriptionsListCmd.Flags().Int("offset", 0, "Offset for pagination")

	// Update flags
	subscriptionsUpdateCmd.Flags().String("plan", "", "New plan ID")
	subscriptionsUpdateCmd.Flags().String("trial-end", "", "Trial end timestamp (Unix timestamp or RFC3339)")
	subscriptionsUpdateCmd.Flags().Bool("prorate", false, "Prorate charges")
	subscriptionsUpdateCmd.Flags().String("metadata", "", "Metadata (key1=value1,key2=value2)")

	// Resume flags
	subscriptionsResumeCmd.Flags().String("trial-end", "", "Trial end timestamp (Unix timestamp or RFC3339)")
	subscriptionsResumeCmd.Flags().Bool("prorate", false, "Prorate charges")
}
