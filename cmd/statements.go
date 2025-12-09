package cmd

import (
	"fmt"

	"github.com/payjp/payjp-cli/internal/client"
	"github.com/payjp/payjp-go/v1"
	"github.com/spf13/cobra"
)

var statementsCmd = &cobra.Command{
	Use:     "statements",
	Aliases: []string{"statement"},
	Short:   "Manage statements",
	Long:    `Retrieve and list transaction statements.`,
}

var statementsGetCmd = &cobra.Command{
	Use:   "get <statement_id>",
	Short: "Get statement information",
	Long: `Retrieve information about a specific statement.

Example:
  payjp statements get st_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		statementID := args[0]

		result, err := client.GetStatement().Retrieve(statementID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var statementsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List statements",
	Long: `List all statements with optional filters.

Example:
  payjp statements list --limit 10
  payjp statements list --owner merchant`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		owner, _ := cmd.Flags().GetString("owner")
		sourceTransfer, _ := cmd.Flags().GetString("source-transfer")

		params := payjp.StatementListParams{}

		if limit > 0 {
			params.Limit = payjp.Int(limit)
		}
		if offset > 0 {
			params.Offset = payjp.Int(offset)
		}
		if owner != "" {
			params.Owner = payjp.String(owner)
		}
		if sourceTransfer != "" {
			params.SourceTransfer = payjp.String(sourceTransfer)
		}

		result, _, err := client.GetStatement().All(&params)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var statementsDownloadUrlCmd = &cobra.Command{
	Use:   "download-url <statement_id>",
	Short: "Get download URL for a statement",
	Long: `Generate a download URL for a specific statement.

Example:
  payjp statements download-url st_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		statementID := args[0]

		statement, err := client.GetStatement().Retrieve(statementID)
		if err != nil {
			handleError(err)
			return nil
		}

		urls, err := statement.StatementUrls()
		if err != nil {
			handleError(err)
			return nil
		}

		if quiet {
			fmt.Println(urls.URL)
			return nil
		}

		return outputResult(map[string]interface{}{
			"id":      statementID,
			"url":     urls.URL,
			"expires": urls.Expires,
		})
	},
}

func init() {
	rootCmd.AddCommand(statementsCmd)

	statementsCmd.AddCommand(statementsGetCmd)
	statementsCmd.AddCommand(statementsListCmd)
	statementsCmd.AddCommand(statementsDownloadUrlCmd)

	// List flags
	statementsListCmd.Flags().Int("limit", 10, "Number of items to return")
	statementsListCmd.Flags().Int("offset", 0, "Offset for pagination")
	statementsListCmd.Flags().String("owner", "", "Filter by owner type (merchant, tenant)")
	statementsListCmd.Flags().String("source-transfer", "", "Filter by source transfer ID")
}
