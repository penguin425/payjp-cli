package cmd

import (
	"fmt"

	"github.com/payjp/payjp-cli/internal/client"
	"github.com/payjp/payjp-cli/internal/util"
	"github.com/payjp/payjp-go/v1"
	"github.com/spf13/cobra"
)

var balancesCmd = &cobra.Command{
	Use:     "balances",
	Aliases: []string{"balance"},
	Short:   "Manage balances",
	Long:    `Retrieve and list account balances.`,
}

var balancesGetCmd = &cobra.Command{
	Use:   "get <balance_id>",
	Short: "Get balance information",
	Long: `Retrieve information about a specific balance.

Example:
  payjp balances get ba_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		balanceID := args[0]

		result, err := client.GetBalance().Retrieve(balanceID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var balancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List balances",
	Long: `List all balances with optional filters.

Example:
  payjp balances list --limit 10
  payjp balances list --owner merchant`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		since, _ := cmd.Flags().GetString("since")
		until, _ := cmd.Flags().GetString("until")
		owner, _ := cmd.Flags().GetString("owner")

		params := payjp.BalanceListParams{}

		if limit > 0 {
			params.Limit = payjp.Int(limit)
		}
		if offset > 0 {
			params.Offset = payjp.Int(offset)
		}
		if since != "" {
			ts, err := util.ParseTimestamp(since)
			if err != nil {
				return err
			}
			params.Since = payjp.Int(int(ts))
		}
		if until != "" {
			ts, err := util.ParseTimestamp(until)
			if err != nil {
				return err
			}
			params.Until = payjp.Int(int(ts))
		}
		if owner != "" {
			params.Owner = payjp.String(owner)
		}

		result, _, err := client.GetBalance().All(&params)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var balancesDownloadUrlCmd = &cobra.Command{
	Use:   "download-url <balance_id>",
	Short: "Get download URL for balance statements",
	Long: `Generate a download URL for all statements in a specific balance.

Example:
  payjp balances download-url ba_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		balanceID := args[0]

		balance, err := client.GetBalance().Retrieve(balanceID)
		if err != nil {
			handleError(err)
			return nil
		}

		urls, err := balance.StatementUrls()
		if err != nil {
			handleError(err)
			return nil
		}

		if quiet {
			fmt.Println(urls.URL)
			return nil
		}

		return outputResult(map[string]interface{}{
			"id":      balanceID,
			"url":     urls.URL,
			"expires": urls.Expires,
		})
	},
}

func init() {
	rootCmd.AddCommand(balancesCmd)

	balancesCmd.AddCommand(balancesGetCmd)
	balancesCmd.AddCommand(balancesListCmd)
	balancesCmd.AddCommand(balancesDownloadUrlCmd)

	// List flags
	balancesListCmd.Flags().Int("limit", 10, "Number of items to return")
	balancesListCmd.Flags().Int("offset", 0, "Offset for pagination")
	balancesListCmd.Flags().String("since", "", "Filter by created timestamp (Unix timestamp or RFC3339)")
	balancesListCmd.Flags().String("until", "", "Filter by created timestamp (Unix timestamp or RFC3339)")
	balancesListCmd.Flags().String("owner", "", "Filter by owner type (merchant, tenant)")
}
