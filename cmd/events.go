package cmd

import (
	"fmt"
	"time"

	"github.com/payjp/payjp-cli/internal/client"
	"github.com/payjp/payjp-cli/internal/util"
	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:     "events",
	Aliases: []string{"event"},
	Short:   "Manage events",
	Long:    `Retrieve and list webhook events.`,
}

var eventsGetCmd = &cobra.Command{
	Use:   "get <event_id>",
	Short: "Get event information",
	Long: `Retrieve information about a specific event.

Example:
  payjp events get evnt_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		eventID := args[0]

		result, err := client.GetEvent().Retrieve(eventID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var eventsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List events",
	Long: `List all events with optional filters.

Example:
  payjp events list --limit 10
  payjp events list --type charge.succeeded
  payjp events list --resource-id ch_xxxxx`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		eventType, _ := cmd.Flags().GetString("type")
		resourceID, _ := cmd.Flags().GetString("resource-id")
		since, _ := cmd.Flags().GetString("since")
		until, _ := cmd.Flags().GetString("until")

		caller := client.GetEvent().List()

		if limit > 0 {
			caller.Limit(limit)
		}
		if offset > 0 {
			caller.Offset(offset)
		}
		if eventType != "" {
			caller.Type(eventType)
		}
		if resourceID != "" {
			caller.ResourceID(resourceID)
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

var eventsTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List available event types",
	Long:  `Display a list of all available event types.`,
	Run: func(cmd *cobra.Command, args []string) {
		eventTypes := []string{
			"charge.succeeded",
			"charge.failed",
			"charge.updated",
			"charge.refunded",
			"charge.captured",
			"customer.created",
			"customer.updated",
			"customer.deleted",
			"customer.card.created",
			"customer.card.updated",
			"customer.card.deleted",
			"plan.created",
			"plan.updated",
			"plan.deleted",
			"subscription.created",
			"subscription.updated",
			"subscription.deleted",
			"subscription.paused",
			"subscription.resumed",
			"subscription.canceled",
			"subscription.renewed",
			"transfer.succeeded",
			"token.created",
		}

		fmt.Println("Available event types:")
		for _, t := range eventTypes {
			fmt.Printf("  %s\n", t)
		}
	},
}

func init() {
	rootCmd.AddCommand(eventsCmd)

	eventsCmd.AddCommand(eventsGetCmd)
	eventsCmd.AddCommand(eventsListCmd)
	eventsCmd.AddCommand(eventsTypesCmd)

	// List flags
	eventsListCmd.Flags().Int("limit", 10, "Number of items to return")
	eventsListCmd.Flags().Int("offset", 0, "Offset for pagination")
	eventsListCmd.Flags().String("type", "", "Filter by event type")
	eventsListCmd.Flags().String("resource-id", "", "Filter by resource ID")
	eventsListCmd.Flags().String("since", "", "Filter by created timestamp (Unix timestamp or RFC3339)")
	eventsListCmd.Flags().String("until", "", "Filter by created timestamp (Unix timestamp or RFC3339)")
}
