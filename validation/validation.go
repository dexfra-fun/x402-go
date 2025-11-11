package validation

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/dexfra-fun/x402-go"
)

var (
	// solanaAddressRegex matches Solana base58 addresses (32-44 chars, base58 charset)
	solanaAddressRegex = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{32,44}$`)
)

// ValidateAmount validates that an amount string is a valid positive integer.
// Returns an error if the amount is empty, malformed, or not greater than zero.
func ValidateAmount(amount string) error {
	if amount == "" {
		return fmt.Errorf("amount cannot be empty")
	}

	// Parse as big.Int to handle large values
	amt := new(big.Int)
	amt, ok := amt.SetString(amount, 10)
	if !ok {
		return fmt.Errorf("invalid amount format: %s", amount)
	}

	if amt.Sign() <= 0 {
		return fmt.Errorf("amount must be greater than 0, got: %s", amount)
	}

	return nil
}

// ValidateAddress validates a Solana address format.
// It uses ValidateNetwork to confirm the network is supported.
func ValidateAddress(address string, network string) error {
	if address == "" {
		return fmt.Errorf("address cannot be empty")
	}

	if err := x402.ValidateNetwork(network); err != nil {
		return fmt.Errorf("cannot validate address: %w", err)
	}

	// Validate based on network (currently only Solana supported)
	if network == "solana" || network == "solana-devnet" {
		if !solanaAddressRegex.MatchString(address) {
			return fmt.Errorf("invalid Solana address format: %s (expected base58 string 32-44 chars)", address)
		}
		return nil
	}

	return fmt.Errorf("unsupported network for address validation: %s", network)
}

// ValidatePaymentRequirement performs comprehensive validation of a payment requirement.
// It validates the amount, network, addresses, scheme, and other required fields.
func ValidatePaymentRequirement(req x402.PaymentRequirement) error {
	// Validate amount
	if err := ValidateAmount(req.MaxAmountRequired); err != nil {
		return fmt.Errorf("invalid requirement: %w", err)
	}

	// Validate network
	if req.Network == "" {
		return fmt.Errorf("invalid requirement: network cannot be empty")
	}

	if err := x402.ValidateNetwork(req.Network); err != nil {
		return fmt.Errorf("invalid requirement: %w", err)
	}

	// Validate recipient address
	if err := ValidateAddress(req.PayTo, req.Network); err != nil {
		return fmt.Errorf("invalid requirement: payTo %w", err)
	}

	// Validate asset address (required)
	if req.Asset == "" {
		return fmt.Errorf("invalid requirement: asset address cannot be empty")
	}

	if err := ValidateAddress(req.Asset, req.Network); err != nil {
		return fmt.Errorf("invalid requirement: asset %w", err)
	}

	// Validate scheme
	switch req.Scheme {
	case "exact", "max", "subscription":
		// Valid schemes
	case "":
		return fmt.Errorf("invalid requirement: scheme cannot be empty")
	default:
		return fmt.Errorf("invalid requirement: unsupported scheme %s", req.Scheme)
	}

	// Validate timeout (must be non-negative)
	if req.MaxTimeoutSeconds < 0 {
		return fmt.Errorf("invalid requirement: timeout cannot be negative: %d", req.MaxTimeoutSeconds)
	}

	return nil
}

// ValidatePaymentPayload validates a payment payload structure.
// It checks the version, scheme, network, and payload fields.
func ValidatePaymentPayload(payment x402.PaymentPayload) error {
	if payment.X402Version != 1 {
		return fmt.Errorf("unsupported x402 version: %d", payment.X402Version)
	}

	if payment.Scheme == "" {
		return fmt.Errorf("scheme cannot be empty")
	}

	if payment.Network == "" {
		return fmt.Errorf("network cannot be empty")
	}

	if err := x402.ValidateNetwork(payment.Network); err != nil {
		return fmt.Errorf("invalid network: %w", err)
	}

	if payment.Payload == nil {
		return fmt.Errorf("payload cannot be nil")
	}

	return nil
}
