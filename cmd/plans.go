package cmd

import (
	"github.com/payjp/payjp-cli/internal/client"
	"github.com/payjp/payjp-cli/internal/util"
	"github.com/payjp/payjp-go/v1"
	"github.com/spf13/cobra"
)

var plansCmd = &cobra.Command{
	Use:     "plans",
	Aliases: []string{"plan"},
	Short:   "Manage subscription plans",
	Long:    `Create, retrieve, update, and delete subscription plans.`,
}

var plansCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new plan",
	Long: `Create a new subscription plan.

Example:
  payjp plans create --amount 1000 --currency jpy --interval month
  payjp plans create --amount 1000 --currency jpy --interval month --name "Basic Plan"
  payjp plans create --amount 1000 --currency jpy --interval month --trial-days 14`,
	RunE: func(cmd *cobra.Command, args []string) error {
		amount, _ := cmd.Flags().GetInt("amount")
		currency, _ := cmd.Flags().GetString("currency")
		interval, _ := cmd.Flags().GetString("interval")
		id, _ := cmd.Flags().GetString("id")
		name, _ := cmd.Flags().GetString("name")
		trialDays, _ := cmd.Flags().GetInt("trial-days")
		billingDay, _ := cmd.Flags().GetInt("billing-day")
		metadata, _ := cmd.Flags().GetString("metadata")

		if err := util.ValidateAmount(amount); err != nil {
			return err
		}
		if err := util.ValidateCurrency(currency); err != nil {
			return err
		}
		if err := util.ValidateInterval(interval); err != nil {
			return err
		}

		plan := payjp.Plan{
			Amount:   amount,
			Currency: currency,
			Interval: interval,
		}

		if id != "" {
			plan.ID = id
		}
		if name != "" {
			plan.Name = name
		}
		if trialDays > 0 {
			plan.TrialDays = trialDays
		}
		if billingDay > 0 {
			plan.BillingDay = billingDay
		}
		if metadata != "" {
			plan.Metadata = util.ParseMetadata(metadata)
		}

		result, err := client.GetPlan().Create(plan)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var plansGetCmd = &cobra.Command{
	Use:   "get <plan_id>",
	Short: "Get plan information",
	Long: `Retrieve information about a specific plan.

Example:
  payjp plans get pln_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		planID := args[0]

		result, err := client.GetPlan().Retrieve(planID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var plansListCmd = &cobra.Command{
	Use:   "list",
	Short: "List plans",
	Long: `List all subscription plans.

Example:
  payjp plans list --limit 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		caller := client.GetPlan().List()

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

var plansUpdateCmd = &cobra.Command{
	Use:   "update <plan_id>",
	Short: "Update plan information",
	Long: `Update information for a specific plan.

Example:
  payjp plans update pln_xxxxx --name "Premium Plan"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		planID := args[0]
		name, _ := cmd.Flags().GetString("name")
		metadata, _ := cmd.Flags().GetString("metadata")

		plan := payjp.Plan{}

		if name != "" {
			plan.Name = name
		}
		if metadata != "" {
			plan.Metadata = util.ParseMetadata(metadata)
		}

		result, err := client.GetPlan().Update(planID, plan)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var plansDeleteCmd = &cobra.Command{
	Use:   "delete <plan_id>",
	Short: "Delete a plan",
	Long: `Delete a specific plan.

Example:
  payjp plans delete pln_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		planID := args[0]

		err := client.GetPlan().Delete(planID)
		if err != nil {
			handleError(err)
			return nil
		}

		if quiet {
			return nil
		}

		return outputResult(map[string]interface{}{
			"id":      planID,
			"deleted": true,
		})
	},
}

func init() {
	rootCmd.AddCommand(plansCmd)

	plansCmd.AddCommand(plansCreateCmd)
	plansCmd.AddCommand(plansGetCmd)
	plansCmd.AddCommand(plansListCmd)
	plansCmd.AddCommand(plansUpdateCmd)
	plansCmd.AddCommand(plansDeleteCmd)

	// Create flags
	plansCreateCmd.Flags().Int("amount", 0, "Amount in smallest currency unit (required)")
	plansCreateCmd.Flags().String("currency", "jpy", "Currency code")
	plansCreateCmd.Flags().String("interval", "month", "Billing interval (month or year)")
	plansCreateCmd.Flags().String("id", "", "Custom plan ID")
	plansCreateCmd.Flags().String("name", "", "Plan name")
	plansCreateCmd.Flags().Int("trial-days", 0, "Trial period in days")
	plansCreateCmd.Flags().Int("billing-day", 0, "Billing day of month (1-31)")
	plansCreateCmd.Flags().String("metadata", "", "Metadata (key1=value1,key2=value2)")
	plansCreateCmd.MarkFlagRequired("amount")

	// List flags
	plansListCmd.Flags().Int("limit", 10, "Number of items to return")
	plansListCmd.Flags().Int("offset", 0, "Offset for pagination")

	// Update flags
	plansUpdateCmd.Flags().String("name", "", "New plan name")
	plansUpdateCmd.Flags().String("metadata", "", "Metadata (key1=value1,key2=value2)")
}
