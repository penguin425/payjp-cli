package cmd

import (
	"fmt"
	"os"

	"github.com/payjp/payjp-cli/internal/client"
	"github.com/payjp/payjp-cli/internal/config"
	"github.com/payjp/payjp-cli/internal/output"
	"github.com/payjp/payjp-cli/internal/util"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time
	Version = "dev"

	// Global flags
	cfgFile   string
	apiKey    string
	outputFmt string
	liveMode  bool
	verbose   bool
	quiet     bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "payjp",
	Short: "PAY.JP CLI - Command line interface for PAY.JP API",
	Long: `PAY.JP CLI is a command line tool for interacting with the PAY.JP payment API.

It allows you to manage charges, customers, cards, plans, subscriptions,
and other PAY.JP resources directly from your terminal.

Example:
  # Set API key
  payjp config set api-key sk_test_xxxxx

  # Create a charge
  payjp charges create --amount 1000 --currency jpy --card tok_xxxxx

  # List customers
  payjp customers list --limit 10

For more information, visit: https://pay.jp/docs/api/`,
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Track if --output flag was explicitly set
		outputFmtChanged = cmd.Flags().Changed("output")

		// Skip client initialization for config commands
		if cmd.Parent() != nil && cmd.Parent().Name() == "config" {
			return nil
		}
		if cmd.Name() == "config" {
			return nil
		}

		// Initialize configuration
		if err := config.Init(cfgFile); err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}

		// Set live mode environment variable if --live flag is used
		if liveMode {
			os.Setenv("PAYJP_LIVE", "true")
		}

		// Initialize client with API key override if provided
		opts := []client.Option{}
		if apiKey != "" {
			opts = append(opts, client.WithAPIKey(apiKey))
		}

		if err := client.Init(opts...); err != nil {
			return err
		}

		return nil
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(int(util.ExitGeneralError))
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ~/.payjp/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "api-key", "k", "", "API key (overrides config file and environment variable)")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "table", "output format (json, table, yaml)")
	rootCmd.PersistentFlags().BoolVar(&liveMode, "live", false, "use live mode (default is test mode)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "quiet output (only output IDs)")
}

func initConfig() {
	// Configuration is initialized in PersistentPreRunE
}

// outputFmtChanged tracks if --output flag was explicitly set
var outputFmtChanged bool

// getOutputFormat returns the output format to use
func getOutputFormat() string {
	if quiet {
		return "quiet"
	}
	if outputFmtChanged {
		return outputFmt
	}
	return config.GetOutputFormat()
}

// outputResult outputs the result in the appropriate format
func outputResult(data interface{}) error {
	format := getOutputFormat()
	return output.Output(format, data)
}

// outputResultQuiet outputs only the ID
func outputResultQuiet(data interface{}) error {
	return output.OutputQuiet(data)
}

// handleError handles errors and exits with appropriate code
func handleError(err error) {
	code := util.HandleError(err)
	fmt.Fprintf(os.Stderr, "\nExit code: %d\n", code)
	os.Exit(int(code))
}

// printVerbose prints verbose output if enabled
func printVerbose(format string, args ...interface{}) {
	if verbose {
		fmt.Printf(format+"\n", args...)
	}
}
