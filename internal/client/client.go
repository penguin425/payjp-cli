package client

import (
	"fmt"

	"github.com/payjp/payjp-cli/internal/config"
	"github.com/payjp/payjp-go/v1"
)

var (
	client *payjp.Service
)

// Options represents client options
type Options struct {
	APIKey       string
	MaxRetry     int
	InitialDelay int
	MaxDelay     int
}

// Option is a function that configures Options
type Option func(*Options)

// WithAPIKey sets the API key
func WithAPIKey(apiKey string) Option {
	return func(o *Options) {
		o.APIKey = apiKey
	}
}

// WithMaxRetry sets the max retry count
func WithMaxRetry(maxRetry int) Option {
	return func(o *Options) {
		o.MaxRetry = maxRetry
	}
}

// WithInitialDelay sets the initial delay for retries
func WithInitialDelay(delay int) Option {
	return func(o *Options) {
		o.InitialDelay = delay
	}
}

// WithMaxDelay sets the max delay for retries
func WithMaxDelay(delay int) Option {
	return func(o *Options) {
		o.MaxDelay = delay
	}
}

// Init initializes the PAY.JP client
func Init(opts ...Option) error {
	retryCfg := config.GetRetryConfig()

	options := &Options{
		APIKey:       config.GetAPIKey(),
		MaxRetry:     retryCfg.MaxCount,
		InitialDelay: retryCfg.InitialDelay,
		MaxDelay:     retryCfg.MaxDelay,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.APIKey == "" {
		return fmt.Errorf("API key is required. Set it via --api-key flag, PAYJP_API_KEY environment variable, or config file")
	}

	client = payjp.New(options.APIKey, nil,
		payjp.WithMaxCount(options.MaxRetry),
		payjp.WithInitialDelay(float64(options.InitialDelay)),
		payjp.WithMaxDelay(float64(options.MaxDelay)),
	)

	return nil
}

// Get returns the PAY.JP client
func Get() *payjp.Service {
	return client
}

// GetCharge returns the Charge service
func GetCharge() *payjp.ChargeService {
	return client.Charge
}

// GetCustomer returns the Customer service
func GetCustomer() *payjp.CustomerService {
	return client.Customer
}

// GetPlan returns the Plan service
func GetPlan() *payjp.PlanService {
	return client.Plan
}

// GetSubscription returns the Subscription service
func GetSubscription() *payjp.SubscriptionService {
	return client.Subscription
}

// GetToken returns the Token service
func GetToken() *payjp.TokenService {
	return client.Token
}

// GetTransfer returns the Transfer service
func GetTransfer() *payjp.TransferService {
	return client.Transfer
}

// GetEvent returns the Event service
func GetEvent() *payjp.EventService {
	return client.Event
}

// GetStatement returns the Statement service
func GetStatement() *payjp.StatementService {
	return client.Statement
}

// GetTerm returns the Term service
func GetTerm() *payjp.TermService {
	return client.Term
}

// GetBalance returns the Balance service
func GetBalance() *payjp.BalanceService {
	return client.Balance
}

// GetAccount returns the Account service
func GetAccount() *payjp.AccountService {
	return client.Account
}
