package x402

import (
	"context"
	"time"

	"github.com/mark3labs/x402-go"
	"github.com/shopspring/decimal"
)

const (
	// defaultCacheTTL is the default cache time-to-live duration.
	defaultCacheTTL = 5 * time.Minute
)

// PricingStrategy defines how to price API resources.
type PricingStrategy interface {
	GetPrice(ctx context.Context, resource Resource) (decimal.Decimal, error)
}

// Resource represents an API endpoint being accessed.
type Resource struct {
	Path   string
	Method string
	Params map[string]string
}

// Config holds the configuration for x402 middleware.
type Config struct {
	// Required fields
	RecipientAddress string
	Network          string
	FacilitatorURL   string
	PricingStrategy  PricingStrategy

	// Optional fields
	CacheTTL time.Duration
	Networks map[string]NetworkConfig
	Logger   Logger
}

// NetworkConfig defines blockchain network configuration.
type NetworkConfig struct {
	ChainID     string
	Name        string
	ChainConfig x402.ChainConfig
}

// Logger defines the logging interface.
type Logger interface {
	Printf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

// DefaultLogger is a no-op logger.
type DefaultLogger struct{}

// Printf is a no-op implementation of the Logger interface.
func (l *DefaultLogger) Printf(_ string, _ ...interface{}) {}

// Errorf is a no-op implementation of the Logger interface.
func (l *DefaultLogger) Errorf(_ string, _ ...interface{}) {}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.RecipientAddress == "" {
		return ErrMissingRecipient
	}
	if c.Network == "" {
		return ErrMissingNetwork
	}
	if c.FacilitatorURL == "" {
		return ErrMissingFacilitator
	}
	if c.PricingStrategy == nil {
		return ErrMissingPricing
	}

	// Set defaults
	if c.CacheTTL == 0 {
		c.CacheTTL = defaultCacheTTL
	}
	if c.Logger == nil {
		c.Logger = &DefaultLogger{}
	}

	return nil
}

// PaymentInfo contains payment metadata.
type PaymentInfo struct {
	Amount    decimal.Decimal
	Currency  string
	Recipient string
	FeePayer  string
}
