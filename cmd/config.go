package cmd

import (
	"fmt"

	"github.com/payjp/payjp-cli/internal/config"
	"github.com/payjp/payjp-cli/internal/util"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  `Manage CLI configuration including API keys, profiles, and output settings.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available keys:
  api-key      Set the API key for the default profile
  output       Set the default output format (json, table, yaml)

Example:
  payjp config set api-key sk_test_xxxxx
  payjp config set output json`,
	Args: cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return config.Init(cfgFile)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		switch key {
		case "api-key":
			profileName := config.Get().DefaultProfile
			if profileName == "" {
				profileName = "default"
			}
			if err := config.SetAPIKey(profileName, value); err != nil {
				return err
			}
			fmt.Printf("API key set for profile '%s'\n", profileName)

		case "output":
			if value != "json" && value != "table" && value != "yaml" {
				return fmt.Errorf("invalid output format: %s (use json, table, or yaml)", value)
			}
			cfg := config.Get()
			cfg.Output.Format = value
			if err := config.Save(); err != nil {
				return err
			}
			fmt.Printf("Output format set to '%s'\n", value)

		default:
			return fmt.Errorf("unknown configuration key: %s", key)
		}

		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current CLI configuration.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return config.Init(cfgFile)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()

		fmt.Println("Configuration:")
		fmt.Println("==============")
		fmt.Printf("Config file: %s\n", config.DefaultConfigPath())
		fmt.Printf("Default profile: %s\n", cfg.DefaultProfile)
		fmt.Printf("Output format: %s\n", cfg.Output.Format)
		fmt.Printf("Color output: %v\n", cfg.Output.Color)
		fmt.Println()

		fmt.Println("Retry settings:")
		fmt.Printf("  Max retries: %d\n", cfg.Retry.MaxCount)
		fmt.Printf("  Initial delay: %ds\n", cfg.Retry.InitialDelay)
		fmt.Printf("  Max delay: %ds\n", cfg.Retry.MaxDelay)
		fmt.Println()

		fmt.Println("Profiles:")
		for name, profile := range cfg.Profiles {
			current := ""
			if name == cfg.DefaultProfile {
				current = " (current)"
			}
			fmt.Printf("  %s%s:\n", name, current)
			fmt.Printf("    API key: %s\n", util.MaskAPIKey(profile.APIKey))
			fmt.Printf("    Mode: %s\n", profile.Mode)
		}

		if len(cfg.Profiles) == 0 {
			fmt.Println("  (none configured)")
		}

		return nil
	},
}

var configSetProfileCmd = &cobra.Command{
	Use:   "set-profile <name>",
	Short: "Create or update a profile",
	Long: `Create or update a named profile with specific settings.

Example:
  payjp config set-profile production --api-key sk_live_xxxxx
  payjp config set-profile development --api-key sk_test_xxxxx`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return config.Init(cfgFile)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		profileAPIKey, _ := cmd.Flags().GetString("api-key")
		mode, _ := cmd.Flags().GetString("mode")

		if profileAPIKey == "" {
			return fmt.Errorf("--api-key is required")
		}

		if mode == "" {
			// Auto-detect mode from API key prefix
			if len(profileAPIKey) > 8 && profileAPIKey[:8] == "sk_live_" {
				mode = "live"
			} else {
				mode = "test"
			}
		} else if mode != "test" && mode != "live" {
			return fmt.Errorf("invalid mode: %s (use 'test' or 'live')", mode)
		}

		profile := config.Profile{
			APIKey: profileAPIKey,
			Mode:   mode,
		}

		if err := config.SetProfile(name, profile); err != nil {
			return err
		}

		fmt.Printf("Profile '%s' saved (mode: %s)\n", name, mode)
		return nil
	},
}

var configUseProfileCmd = &cobra.Command{
	Use:   "use-profile <name>",
	Short: "Switch to a profile",
	Long: `Set the default profile to use.

Example:
  payjp config use-profile production`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return config.Init(cfgFile)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if err := config.UseProfile(name); err != nil {
			return err
		}

		fmt.Printf("Now using profile '%s'\n", name)
		return nil
	},
}

var configListProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: "List all profiles",
	Long:  `Display a list of all configured profiles.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return config.Init(cfgFile)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		profiles := config.ListProfiles()

		if len(profiles) == 0 {
			fmt.Println("No profiles configured.")
			fmt.Println("Use 'payjp config set-profile <name> --api-key <key>' to create one.")
			return nil
		}

		fmt.Println("Profiles:")
		for _, name := range profiles {
			profile := cfg.Profiles[name]
			current := ""
			if name == cfg.DefaultProfile {
				current = " *"
			}
			fmt.Printf("  %s%s (%s)\n", name, current, profile.Mode)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetProfileCmd)
	configCmd.AddCommand(configUseProfileCmd)
	configCmd.AddCommand(configListProfilesCmd)

	// Flags for set-profile
	configSetProfileCmd.Flags().String("api-key", "", "API key for the profile")
	configSetProfileCmd.Flags().String("mode", "", "Mode (test or live, auto-detected from key if not specified)")
}
