// Package x402 provides helper functions and constants for configuring x402 payments
// with USDC across multiple blockchain networks. This package simplifies client and
// middleware setup by providing verified USDC addresses, EIP-3009 parameters, and
// utility functions for creating payment requirements and token configurations.
package x402

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

const (
	// USDCDecimals is the standard number of decimals for USDC tokens.
	USDCDecimals = 6

	// DefaultMaxTimeoutSeconds is the default maximum timeout for payment authorization in seconds.
	DefaultMaxTimeoutSeconds = 300

	// DefaultMimeType is the default MIME type for payment responses.
	DefaultMimeType = "application/json"

	// DefaultScheme is the default payment scheme.
	DefaultScheme = "exact"
)

// ChainConfig contains chain-specific configuration for USDC tokens and payment requirements.
// All USDC addresses and EIP-3009 parameters were verified on 2025-10-28.
type ChainConfig struct {
	// NetworkID is the x402 protocol network identifier (e.g., "base", "solana").
	NetworkID string

	// USDCAddress is the official Circle USDC contract address or mint address.
	USDCAddress string

	// Decimals is the number of decimal places for USDC (always 6).
	Decimals uint8

	// EIP3009Name is the EIP-3009 domain parameter "name" (empty for non-EVM chains).
	EIP3009Name string

	// EIP3009Version is the EIP-3009 domain parameter "version" (empty for non-EVM chains).
	EIP3009Version string
}

// USDCRequirementConfig is the configuration for creating a USDC PaymentRequirement.
// This is a convenience helper for USDC payments. For other tokens, construct
// PaymentRequirement directly.
type USDCRequirementConfig struct {
	// Chain is the chain configuration with USDC details (required).
	Chain ChainConfig

	// Amount is the human-readable USDC amount (e.g., "1.5" = 1.5 USDC).
	// Zero amounts ("0" or "0.0") are allowed for free-with-signature authorization flows.
	Amount string

	// RecipientAddress is the payment recipient address (required).
	RecipientAddress string

	// Description is a human-readable description of the payment (optional).
	Description string

	// Scheme is the payment scheme (optional, defaults to "exact").
	Scheme string

	// MaxTimeoutSeconds is the maximum payment timeout (optional, defaults to 300).
	MaxTimeoutSeconds uint32

	// MimeType is the response MIME type (optional, defaults to "application/json").
	MimeType string
}

// Solana chain configurations.
var (
	// SolanaMainnet is the configuration for Solana mainnet.
	// USDC address verified 2025-10-28.
	SolanaMainnet = ChainConfig{
		NetworkID:      "solana",
		USDCAddress:    "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Decimals:       USDCDecimals,
		EIP3009Name:    "",
		EIP3009Version: "",
	}

	// SolanaDevnet is the configuration for Solana devnet.
	// USDC address verified 2025-10-28.
	SolanaDevnet = ChainConfig{
		NetworkID:      "solana-devnet",
		USDCAddress:    "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU",
		Decimals:       USDCDecimals,
		EIP3009Name:    "",
		EIP3009Version: "",
	}
)

// NewUSDCTokenConfig creates a TokenConfig for USDC on the given chain with the specified priority.
// This is a convenience helper for USDC. For other tokens, construct TokenConfig directly.
// The returned TokenConfig has:
//   - Address set to the chain's USDC address
//   - Symbol set to "USDC"
//   - Decimals set to 6
//   - Priority set to the provided value (lower numbers = higher priority)
func NewUSDCTokenConfig(chain ChainConfig, priority int) TokenConfig {
	return TokenConfig{
		Address:  chain.USDCAddress,
		Symbol:   "USDC",
		Decimals: USDCDecimals,
		Priority: priority,
	}
}

// NewUSDCPaymentRequirement creates a PaymentRequirement for USDC from the given configuration.
// This is a convenience helper for USDC payments. For other tokens, construct PaymentRequirement directly.
// It validates inputs, converts the amount to atomic units (assuming 6 decimals for USDC),
// applies defaults for optional fields, and populates EIP-3009 parameters for EVM chains.
//
// Amount conversion uses standard float64 rounding (banker's rounding) for precision beyond 6 decimals.
// Zero amounts ("0" or "0.0") are explicitly allowed for free-with-signature authorization flows.
//
// Default values:
//   - Scheme: "exact"
//   - MaxTimeoutSeconds: 300
//   - MimeType: "application/json"
//
// Returns an error if validation fails. Error format: "parameterName: reason".
func NewUSDCPaymentRequirement(config USDCRequirementConfig) (PaymentRequirement, error) {
	// Validate recipient address
	if config.RecipientAddress == "" {
		return PaymentRequirement{}, errors.New("recipientAddress: cannot be empty")
	}

	// Parse and validate amount
	amount, err := strconv.ParseFloat(config.Amount, 64)
	if err != nil {
		return PaymentRequirement{}, errors.New("amount: invalid format")
	}
	if amount < 0 {
		return PaymentRequirement{}, errors.New("amount: must be non-negative")
	}

	// Convert to atomic units (USDC always has 6 decimals)
	const usdcMultiplier = 1e6
	atomicUnits := uint64(math.RoundToEven(amount * usdcMultiplier))
	atomicString := strconv.FormatUint(atomicUnits, 10)

	// Apply defaults
	scheme := config.Scheme
	if scheme == "" {
		scheme = DefaultScheme
	}

	maxTimeout := config.MaxTimeoutSeconds
	if maxTimeout == 0 {
		maxTimeout = DefaultMaxTimeoutSeconds
	}

	mimeType := config.MimeType
	if mimeType == "" {
		mimeType = DefaultMimeType
	}

	// Create base payment requirement
	req := PaymentRequirement{
		Scheme:            scheme,
		Network:           config.Chain.NetworkID,
		MaxAmountRequired: atomicString,
		Asset:             config.Chain.USDCAddress,
		PayTo:             config.RecipientAddress,
		Description:       config.Description,
		MimeType:          mimeType,
		MaxTimeoutSeconds: int(maxTimeout),
	}

	// Populate EIP-3009 extra field for EVM chains
	if config.Chain.EIP3009Name != "" {
		req.Extra = map[string]any{
			"name":    config.Chain.EIP3009Name,
			"version": config.Chain.EIP3009Version,
		}
	}

	return req, nil
}

// ValidateNetwork validates a network identifier.
// Returns nil if the network is supported, error otherwise.
//
// Supported networks:
//   - solana
//   - solana-devnet
func ValidateNetwork(networkID string) error {
	if networkID == "" {
		return errors.New("networkID: cannot be empty")
	}

	// Supported networks
	supportedNetworks := map[string]bool{
		"solana":        true,
		"solana-devnet": true,
	}

	if !supportedNetworks[networkID] {
		return fmt.Errorf("networkID: unsupported network '%s'", networkID)
	}

	return nil
}
