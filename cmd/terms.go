package cmd

import (
	"github.com/payjp/payjp-cli/internal/client"
	"github.com/payjp/payjp-go/v1"
	"github.com/spf13/cobra"
)

var termsCmd = &cobra.Command{
	Use:     "terms",
	Aliases: []string{"term"},
	Short:   "Manage terms",
	Long:    `Retrieve and list aggregation terms (billing periods).`,
}

var termsGetCmd = &cobra.Command{
	Use:   "get <term_id>",
	Short: "Get term information",
	Long: `Retrieve information about a specific term.

Example:
  payjp terms get tm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		termID := args[0]

		result, err := client.GetTerm().Retrieve(termID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

var termsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List terms",
	Long: `List all terms with optional filters.

Example:
  payjp terms list --limit 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		params := payjp.TermListParams{}

		if limit > 0 {
			params.Limit = payjp.Int(limit)
		}
		if offset > 0 {
			params.Offset = payjp.Int(offset)
		}

		result, _, err := client.GetTerm().All(&params)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

func init() {
	rootCmd.AddCommand(termsCmd)

	termsCmd.AddCommand(termsGetCmd)
	termsCmd.AddCommand(termsListCmd)

	// List flags
	termsListCmd.Flags().Int("limit", 10, "Number of items to return")
	termsListCmd.Flags().Int("offset", 0, "Offset for pagination")
}
