package cmd

import (
	"github.com/payjp/payjp-cli/internal/client"
	"github.com/spf13/cobra"
)

var accountsCmd = &cobra.Command{
	Use:     "accounts",
	Aliases: []string{"account"},
	Short:   "Manage account",
	Long:    `Retrieve account information.`,
}

var accountsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get account information",
	Long: `Retrieve information about your account.

Example:
  payjp accounts get`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.GetAccount().Retrieve()
		if err != nil {
			handleError(err)
			return nil
		}

		return outputResult(result)
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)

	accountsCmd.AddCommand(accountsGetCmd)
}
