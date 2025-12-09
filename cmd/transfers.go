package cmd

import (
	"time"

	"github.com/payjp/payjp-cli/internal/client"
	"github.com/payjp/payjp-cli/internal/util"
	"github.com/spf13/cobra"
)

var transfersCmd = &cobra.Command{
	Use:     "transfers",
	Aliases: []string{"transfer"},
	Short:   "Manage transfers",
	Long:    `Retrieve and list transfers (payouts to your bank account).`,
}

var transfersGetCmd = &cobra.Command{
	Use:   "get <transfer_id>",
	Short: "Get transfer information",
	Long: `Retrieve information about a specific transfer.

Example:
  payjp transfers get tr_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		transferID := args[0]

		result, err := client.GetTransfer().Retrieve(transferID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var transfersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List transfers",
	Long: `List all transfers with optional filters.

Example:
  payjp transfers list --limit 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		since, _ := cmd.Flags().GetString("since")
		until, _ := cmd.Flags().GetString("until")

		caller := client.GetTransfer().List()

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

func init() {
	rootCmd.AddCommand(transfersCmd)

	transfersCmd.AddCommand(transfersGetCmd)
	transfersCmd.AddCommand(transfersListCmd)

	// List flags
	transfersListCmd.Flags().Int("limit", 10, "Number of items to return")
	transfersListCmd.Flags().Int("offset", 0, "Offset for pagination")
	transfersListCmd.Flags().String("since", "", "Filter by created timestamp (Unix timestamp or RFC3339)")
	transfersListCmd.Flags().String("until", "", "Filter by created timestamp (Unix timestamp or RFC3339)")
}
