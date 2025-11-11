package validation

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"

	"github.com/dexfra-fun/x402-go"
)

var (
	// solanaAddressRegex matches Solana base58 addresses (32-44 chars, base58 charset).
	solanaAddressRegex = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{32,44}$`)
)

const (
	// DecimalBase is the base for parsing big integers.
	DecimalBase = 10
)

// ValidateAmount validates that an amount string is a valid positive integer.
// Returns an error if the amount is empty, malformed, or not greater than zero.
func ValidateAmount(amount string) error {
	if amount == "" {
		return errors.New("amount cannot be empty")
	}

	// Parse as big.Int to handle large values
	amt := new(big.Int)
	amt, ok := amt.SetString(amount, DecimalBase)
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
		return errors.New("address cannot be empty")
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

// validateRequirementNetwork validates the network field of a payment requirement.
func validateRequirementNetwork(network string) error {
	if network == "" {
		return errors.New("invalid requirement: network cannot be empty")
	}
	if err := x402.ValidateNetwork(network); err != nil {
		return fmt.Errorf("invalid requirement: %w", err)
	}
	return nil
}

// validateRequirementAddresses validates the addresses in a payment requirement.
func validateRequirementAddresses(payTo, asset, network string) error {
	if err := ValidateAddress(payTo, network); err != nil {
		return fmt.Errorf("invalid requirement: payTo %w", err)
	}

	if asset == "" {
		return errors.New("invalid requirement: asset address cannot be empty")
	}

	if err := ValidateAddress(asset, network); err != nil {
		return fmt.Errorf("invalid requirement: asset %w", err)
	}
	return nil
}

// validateRequirementScheme validates the scheme field of a payment requirement.
func validateRequirementScheme(scheme string) error {
	switch scheme {
	case "exact", "max", "subscription":
		return nil
	case "":
		return errors.New("invalid requirement: scheme cannot be empty")
	default:
		return fmt.Errorf("invalid requirement: unsupported scheme %s", scheme)
	}
}

// ValidatePaymentRequirement performs comprehensive validation of a payment requirement.
// It validates the amount, network, addresses, scheme, and other required fields.
func ValidatePaymentRequirement(req x402.PaymentRequirement) error {
	// Validate amount
	if err := ValidateAmount(req.MaxAmountRequired); err != nil {
		return fmt.Errorf("invalid requirement: %w", err)
	}

	// Validate network
	if err := validateRequirementNetwork(req.Network); err != nil {
		return err
	}

	// Validate addresses
	if err := validateRequirementAddresses(req.PayTo, req.Asset, req.Network); err != nil {
		return err
	}

	// Validate scheme
	if err := validateRequirementScheme(req.Scheme); err != nil {
		return err
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
		return errors.New("scheme cannot be empty")
	}

	if payment.Network == "" {
		return errors.New("network cannot be empty")
	}

	if err := x402.ValidateNetwork(payment.Network); err != nil {
		return fmt.Errorf("invalid network: %w", err)
	}

	if payment.Payload == nil {
		return errors.New("payload cannot be nil")
	}

	return nil
}
