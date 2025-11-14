package x402

import (
	"context"
	"errors"
	"fmt"
	"strings"

	x402 "github.com/dexfra-fun/x402-go"
	"github.com/mr-tron/base58"
	"github.com/shopspring/decimal"
)

// Middleware handles x402 payment verification.
type Middleware struct {
	config      *Config
	facilitator *FacilitatorClient
	cache       *FeePayerCache
	chainConfig x402.ChainConfig
}

// New creates a new x402 middleware instance.
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
	const cacheCleanupDivisor = 2
	cache.StartCleanupRoutine(config.CacheTTL / cacheCleanupDivisor)

	// Create facilitator client
	facilitator := NewFacilitatorClient(config.FacilitatorURL, cache, config.Logger)

	return &Middleware{
		config:      config,
		facilitator: facilitator,
		cache:       cache,
		chainConfig: chainConfig,
	}, nil
}

// getFeePayer retrieves the fee payer, trying facilitator first then config fallback.
func (m *Middleware) getFeePayer(ctx context.Context) (string, error) {
	// Try facilitator first
	feePayer, err := m.facilitator.GetFeePayer(ctx, m.config.Network)
	if err != nil && !errors.Is(err, ErrFeePayerNotFound) {
		return "", fmt.Errorf("get fee payer: %w", err)
	}

	// Use fallback if facilitator doesn't have one
	if errors.Is(err, ErrFeePayerNotFound) && m.config.FeePayer != "" {
		m.config.Logger.Printf("[x402] Facilitator doesn't provide fee payer, using config fallback")
		return m.config.FeePayer, nil
	}

	return feePayer, nil
}

// validateFeePayer validates the fee payer address.
func (m *Middleware) validateFeePayer(feePayer string) error {
	feePayer = strings.TrimSpace(feePayer)
	if feePayer == "" {
		return ErrMissingFeePayer
	}
	if _, err := base58.Decode(feePayer); err != nil {
		m.config.Logger.Errorf("[x402] Invalid fee payer (not base58): %q: %v", feePayer, err)
		return ErrInvalidFeePayer
	}
	return nil
}

// addRequirementMetadata adds schema and extra metadata to requirement.
func (m *Middleware) addRequirementMetadata(
	ctx context.Context,
	requirement *x402.PaymentRequirement,
	resource Resource,
	feePayer string,
) {
	// Add fee payer to extra metadata
	if requirement.Extra == nil {
		requirement.Extra = make(map[string]any)
	}
	requirement.Extra["feePayer"] = feePayer

	// Add schema if SchemaProvider is configured
	if m.config.SchemaProvider != nil {
		schema, err := m.config.SchemaProvider.GetSchema(ctx, resource)
		if err != nil {
			m.config.Logger.Printf("[x402] Failed to get schema: %v", err)
		} else if schema != nil {
			requirement.OutputSchema = schema
			m.config.Logger.Printf("[x402] Schema added to payment requirement")
		}
	}
}

// ProcessRequest handles payment requirement for a resource.
func (m *Middleware) ProcessRequest(
	ctx context.Context,
	resource Resource,
) (*x402.PaymentRequirement, *PaymentInfo, error) {
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

	// Get and validate fee payer
	feePayer, err := m.getFeePayer(ctx)
	if err != nil {
		return nil, nil, err
	}

	feePayer = strings.TrimSpace(feePayer)
	if err := m.validateFeePayer(feePayer); err != nil {
		return nil, nil, err
	}

	// Create USDC payment requirement
	requirement, err := x402.NewUSDCPaymentRequirement(x402.USDCRequirementConfig{
		Chain:            m.chainConfig,
		Amount:           price.String(),
		RecipientAddress: m.config.RecipientAddress,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create payment requirement: %w", err)
	}

	// Add metadata (feePayer in extra, schema if configured)
	m.addRequirementMetadata(ctx, &requirement, resource, feePayer)

	paymentInfo := &PaymentInfo{
		Amount:    price,
		Currency:  "USDC",
		Recipient: m.config.RecipientAddress,
		FeePayer:  feePayer,
	}

	return &requirement, paymentInfo, nil
}

// GetConfig returns the middleware configuration.
func (m *Middleware) GetConfig() *Config {
	return m.config
}

// GetFacilitator returns the facilitator client.
func (m *Middleware) GetFacilitator() *FacilitatorClient {
	return m.facilitator
}
