package x402

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/x402-go"
	"github.com/mr-tron/base58"
	"github.com/shopspring/decimal"
)

// Middleware handles x402 payment verification
type Middleware struct {
	config      *Config
	facilitator *FacilitatorClient
	cache       *FeePayerCache
	chainConfig x402.ChainConfig
}

// New creates a new x402 middleware instance
func New(config *Config) (*Middleware, error) {
	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Map network to chain config
	chainConfig, err := MapNetworkToChain(config.Network)
	if err != nil {
		return nil, err
	}

	// Create cache
	cache := NewFeePayerCache(config.CacheTTL)
	
	// Start cleanup routine (runs every CacheTTL/2)
	cache.StartCleanupRoutine(config.CacheTTL / 2)

	// Create facilitator client
	facilitator := NewFacilitatorClient(config.FacilitatorURL, cache, config.Logger)

	return &Middleware{
		config:      config,
		facilitator: facilitator,
		cache:       cache,
		chainConfig: chainConfig,
	}, nil
}

// ProcessRequest handles payment requirement for a resource
func (m *Middleware) ProcessRequest(ctx context.Context, resource Resource) (*x402.PaymentRequirement, *PaymentInfo, error) {
	// Get price for this resource
	price, err := m.config.PricingStrategy.GetPrice(ctx, resource)
	if err != nil {
		return nil, nil, fmt.Errorf("get price: %w", err)
	}

	// Free endpoint - no payment required
	if price.LessThanOrEqual(decimal.Zero) {
		return nil, nil, nil
	}

	m.config.Logger.Printf("[x402] Payment required: path=%s method=%s price=%s USDC",
		resource.Path, resource.Method, price.String())

	// Get fee payer from facilitator
	feePayer, err := m.facilitator.GetFeePayer(ctx, m.config.Network)
	if err != nil {
		return nil, nil, fmt.Errorf("get fee payer: %w", err)
	}

	// Validate fee payer address
	feePayer = strings.TrimSpace(feePayer)
	if _, err := base58.Decode(feePayer); err != nil {
		m.config.Logger.Errorf("[x402] Invalid fee payer (not base58): %q: %v", feePayer, err)
		return nil, nil, ErrInvalidFeePayer
	}

	// Convert price to string
	amountStr := price.String()

	// Create USDC payment requirement
	requirement, err := x402.NewUSDCPaymentRequirement(x402.USDCRequirementConfig{
		Chain:            m.chainConfig,
		Amount:           amountStr,
		RecipientAddress: m.config.RecipientAddress,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create payment requirement: %w", err)
	}

	// Add fee payer to extra metadata
	if requirement.Extra == nil {
		requirement.Extra = make(map[string]interface{})
	}
	requirement.Extra["feePayer"] = feePayer

	paymentInfo := &PaymentInfo{
		Amount:    price,
		Currency:  "USDC",
		Recipient: m.config.RecipientAddress,
		FeePayer:  feePayer,
	}

	return &requirement, paymentInfo, nil
}

// GetConfig returns the middleware configuration
func (m *Middleware) GetConfig() *Config {
	return m.config
}

// GetFacilitator returns the facilitator client
func (m *Middleware) GetFacilitator() *FacilitatorClient {
	return m.facilitator
}
