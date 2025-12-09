package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the CLI configuration
type Config struct {
	DefaultProfile string            `mapstructure:"default_profile"`
	Output         OutputConfig      `mapstructure:"output"`
	Retry          RetryConfig       `mapstructure:"retry"`
	Profiles       map[string]Profile `mapstructure:"profiles"`
	Aliases        map[string]string  `mapstructure:"aliases"`
}

// OutputConfig represents output settings
type OutputConfig struct {
	Format string `mapstructure:"format"`
	Color  bool   `mapstructure:"color"`
}

// RetryConfig represents retry settings
type RetryConfig struct {
	MaxCount     int `mapstructure:"max_count"`
	InitialDelay int `mapstructure:"initial_delay"`
	MaxDelay     int `mapstructure:"max_delay"`
}

// Profile represents an API profile
type Profile struct {
	APIKey string `mapstructure:"api_key"`
	Mode   string `mapstructure:"mode"`
}

var (
	cfg        *Config
	configPath string
)

// DefaultConfigDir returns the default configuration directory
func DefaultConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".payjp"
	}
	return filepath.Join(home, ".payjp")
}

// DefaultConfigPath returns the default configuration file path
func DefaultConfigPath() string {
	return filepath.Join(DefaultConfigDir(), "config.yaml")
}

// Init initializes the configuration
func Init(cfgFile string) error {
	if cfgFile != "" {
		configPath = cfgFile
		viper.SetConfigFile(cfgFile)
	} else {
		// Check environment variable
		if envConfig := os.Getenv("PAYJP_CONFIG"); envConfig != "" {
			configPath = envConfig
			viper.SetConfigFile(envConfig)
		} else {
			configPath = DefaultConfigPath()
			viper.AddConfigPath(DefaultConfigDir())
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
		}
	}

	// Set default values
	viper.SetDefault("default_profile", "default")
	viper.SetDefault("output.format", "table")
	viper.SetDefault("output.color", true)
	viper.SetDefault("retry.max_count", 3)
	viper.SetDefault("retry.initial_delay", 2)
	viper.SetDefault("retry.max_delay", 32)

	// Read environment variables
	viper.SetEnvPrefix("PAYJP")
	viper.AutomaticEnv()

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}

// Get returns the current configuration
func Get() *Config {
	if cfg == nil {
		cfg = &Config{
			DefaultProfile: "default",
			Output: OutputConfig{
				Format: "table",
				Color:  true,
			},
			Retry: RetryConfig{
				MaxCount:     3,
				InitialDelay: 2,
				MaxDelay:     32,
			},
			Profiles: make(map[string]Profile),
			Aliases:  make(map[string]string),
		}
	}
	return cfg
}

// GetAPIKey returns the API key to use
func GetAPIKey() string {
	// Priority: environment variable > profile
	if apiKey := os.Getenv("PAYJP_API_KEY"); apiKey != "" {
		return apiKey
	}

	cfg := Get()
	profileName := os.Getenv("PAYJP_PROFILE")
	if profileName == "" {
		profileName = cfg.DefaultProfile
	}

	if profile, ok := cfg.Profiles[profileName]; ok {
		return profile.APIKey
	}

	return ""
}

// GetCurrentProfile returns the current profile
func GetCurrentProfile() (string, *Profile) {
	cfg := Get()
	profileName := os.Getenv("PAYJP_PROFILE")
	if profileName == "" {
		profileName = cfg.DefaultProfile
	}

	if profile, ok := cfg.Profiles[profileName]; ok {
		return profileName, &profile
	}

	return profileName, nil
}

// SetAPIKey sets the API key for the specified profile
func SetAPIKey(profileName, apiKey string) error {
	cfg := Get()
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}

	profile := cfg.Profiles[profileName]
	profile.APIKey = apiKey
	cfg.Profiles[profileName] = profile

	return Save()
}

// SetProfile creates or updates a profile
func SetProfile(name string, profile Profile) error {
	cfg := Get()
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}

	cfg.Profiles[name] = profile
	return Save()
}

// UseProfile sets the default profile
func UseProfile(name string) error {
	cfg := Get()
	if _, ok := cfg.Profiles[name]; !ok {
		return fmt.Errorf("profile '%s' not found", name)
	}

	cfg.DefaultProfile = name
	return Save()
}

// ListProfiles returns all profile names
func ListProfiles() []string {
	cfg := Get()
	profiles := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		profiles = append(profiles, name)
	}
	return profiles
}

// Save saves the configuration to file
func Save() error {
	cfg := Get()

	// Ensure config directory exists with secure permissions
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	viper.Set("default_profile", cfg.DefaultProfile)
	viper.Set("output", cfg.Output)
	viper.Set("retry", cfg.Retry)
	viper.Set("profiles", cfg.Profiles)
	viper.Set("aliases", cfg.Aliases)

	// Write to a temp file first with secure permissions, then rename
	// This prevents a race condition where the file is readable before chmod
	tempFile := configPath + ".tmp"

	// Create temp file with secure permissions (0600) from the start
	f, err := os.OpenFile(tempFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error creating temp config file: %w", err)
	}
	f.Close()

	if err := viper.WriteConfigAs(tempFile); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("error writing config file: %w", err)
	}

	// Ensure temp file has correct permissions (viper may have changed them)
	if err := os.Chmod(tempFile, 0600); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("error setting config file permissions: %w", err)
	}

	// Atomically rename temp file to final path
	if err := os.Rename(tempFile, configPath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("error renaming config file: %w", err)
	}

	return nil
}

// GetOutputFormat returns the output format
func GetOutputFormat() string {
	if format := os.Getenv("PAYJP_OUTPUT"); format != "" {
		return format
	}
	return Get().Output.Format
}

// IsLiveMode returns true if live mode is enabled
func IsLiveMode() bool {
	if live := os.Getenv("PAYJP_LIVE"); live == "true" {
		return true
	}
	_, profile := GetCurrentProfile()
	if profile != nil {
		return profile.Mode == "live"
	}
	return false
}

// GetRetryConfig returns the retry configuration
func GetRetryConfig() RetryConfig {
	return Get().Retry
}

// ResolveAlias resolves a command alias
func ResolveAlias(cmd string) string {
	cfg := Get()
	if resolved, ok := cfg.Aliases[cmd]; ok {
		return resolved
	}
	return cmd
}
