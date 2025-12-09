package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/payjp/payjp-go/v1"
)

// ExitCode represents CLI exit codes
type ExitCode int

const (
	ExitSuccess          ExitCode = 0
	ExitGeneralError     ExitCode = 1
	ExitArgumentError    ExitCode = 2
	ExitConfigError      ExitCode = 3
	ExitAuthError        ExitCode = 4  // 401
	ExitRequestError     ExitCode = 5  // 400
	ExitPaymentError     ExitCode = 6  // 402
	ExitNotFoundError    ExitCode = 7  // 404
	ExitRateLimitError   ExitCode = 8  // 429
	ExitServerError      ExitCode = 9  // 500
)

// Exit exits the program with the given code
func Exit(code ExitCode) {
	os.Exit(int(code))
}

// HandleError handles API errors and returns the appropriate exit code
func HandleError(err error) ExitCode {
	if err == nil {
		return ExitSuccess
	}

	// Check if it's a PAY.JP error
	if payjpErr, ok := err.(*payjp.Error); ok {
		fmt.Fprintf(os.Stderr, "Error: %s\n", payjpErr.Message)
		fmt.Fprintf(os.Stderr, "  Status: %d\n", payjpErr.Status)
		fmt.Fprintf(os.Stderr, "  Type: %s\n", payjpErr.Type)
		if payjpErr.Code != "" {
			fmt.Fprintf(os.Stderr, "  Code: %s\n", payjpErr.Code)
		}
		if payjpErr.Param != "" {
			fmt.Fprintf(os.Stderr, "  Param: %s\n", payjpErr.Param)
		}

		switch payjpErr.Status {
		case 400:
			return ExitRequestError
		case 401:
			return ExitAuthError
		case 402:
			return ExitPaymentError
		case 404:
			return ExitNotFoundError
		case 429:
			return ExitRateLimitError
		default:
			if payjpErr.Status >= 500 {
				return ExitServerError
			}
			return ExitGeneralError
		}
	}

	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	return ExitGeneralError
}

// ParseMetadata parses a metadata string into a map
// Format: key1=value1,key2=value2
func ParseMetadata(s string) map[string]string {
	if s == "" {
		return nil
	}

	metadata := make(map[string]string)
	pairs := strings.Split(s, ",")

	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" {
				metadata[key] = value
			}
		}
	}

	if len(metadata) == 0 {
		return nil
	}

	return metadata
}

// ParseTimestamp parses a timestamp string
// Accepts Unix timestamp or RFC3339 format
func ParseTimestamp(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}

	// Try Unix timestamp first
	if ts, err := strconv.ParseInt(s, 10, 64); err == nil {
		return ts, nil
	}

	// Try RFC3339 format
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return 0, fmt.Errorf("invalid timestamp format: %s (use Unix timestamp or RFC3339)", s)
	}

	return t.Unix(), nil
}

// FormatTimestamp formats a Unix timestamp as a string
func FormatTimestamp(ts int64) string {
	if ts == 0 {
		return ""
	}
	t := time.Unix(ts, 0)
	return t.Format("2006-01-02 15:04:05")
}

// FormatAmount formats an amount with currency
func FormatAmount(amount int, currency string) string {
	switch strings.ToLower(currency) {
	case "jpy":
		return fmt.Sprintf("¥%d", amount)
	case "usd":
		return fmt.Sprintf("$%.2f", float64(amount)/100)
	default:
		return fmt.Sprintf("%d %s", amount, strings.ToUpper(currency))
	}
}

// BoolPtr returns a pointer to a bool
func BoolPtr(b bool) *bool {
	return &b
}

// StringPtr returns a pointer to a string
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// IntPtr returns a pointer to an int
func IntPtr(i int) *int {
	return &i
}

// Int64Ptr returns a pointer to an int64
func Int64Ptr(i int64) *int64 {
	return &i
}

// MaskAPIKey masks an API key for display
func MaskAPIKey(key string) string {
	if len(key) < 8 {
		return "****"
	}
	return key[:7] + "****" + key[len(key)-4:]
}

// ConfirmAction prompts for confirmation
func ConfirmAction(message string) bool {
	fmt.Printf("%s [y/N]: ", message)
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// TruncateString truncates a string to the specified rune length
func TruncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}

// ValidateAmount validates an amount
func ValidateAmount(amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	return nil
}

// ValidateCurrency validates a currency code
func ValidateCurrency(currency string) error {
	currency = strings.ToLower(currency)
	validCurrencies := []string{"jpy", "usd"}
	for _, c := range validCurrencies {
		if currency == c {
			return nil
		}
	}
	return fmt.Errorf("invalid currency: %s (supported: jpy, usd)", currency)
}

// ValidateInterval validates a subscription interval
func ValidateInterval(interval string) error {
	validIntervals := []string{"month", "year"}
	for _, i := range validIntervals {
		if interval == i {
			return nil
		}
	}
	return fmt.Errorf("invalid interval: %s (supported: month, year)", interval)
}
