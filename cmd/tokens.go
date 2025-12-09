package cmd

import (
	"github.com/payjp/payjp-cli/internal/client"
	"github.com/spf13/cobra"
)

var tokensCmd = &cobra.Command{
	Use:     "tokens",
	Aliases: []string{"token"},
	Short:   "Manage tokens",
	Long: `Retrieve token information.

Note: Token creation should be done client-side using PAY.JP Checkout or the JavaScript library.`,
}

var tokensGetCmd = &cobra.Command{
	Use:   "get <token_id>",
	Short: "Get token information",
	Long: `Retrieve information about a specific token.

Example:
  payjp tokens get tok_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenID := args[0]

		result, err := client.GetToken().Retrieve(tokenID)
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

func init() {
	rootCmd.AddCommand(tokensCmd)

	tokensCmd.AddCommand(tokensGetCmd)
}
