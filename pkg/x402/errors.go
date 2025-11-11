package x402

import "errors"

var (
	// ErrMissingRecipient indicates that the recipient address is not configured.
	ErrMissingRecipient = errors.New("x402: recipient address is required")
	// ErrMissingNetwork indicates that the network is not configured.
	ErrMissingNetwork = errors.New("x402: network is required")
	// ErrMissingFacilitator indicates that the facilitator URL is not configured.
	ErrMissingFacilitator = errors.New("x402: facilitator URL is required")
	// ErrMissingPricing indicates that the pricing strategy is not configured.
	ErrMissingPricing = errors.New("x402: pricing strategy is required")

	// ErrInvalidFeePayer indicates that the fee payer address is invalid.
	ErrInvalidFeePayer = errors.New("x402: invalid fee payer address")
	// ErrFacilitatorUnavailable indicates that the facilitator service is unavailable.
	ErrFacilitatorUnavailable = errors.New("x402: facilitator service unavailable")
	// ErrFeePayerNotFound indicates that no fee payer was found for the specified network.
	ErrFeePayerNotFound = errors.New("x402: fee payer not found for network")
	// ErrPaymentVerificationFailed indicates that payment verification failed.
	ErrPaymentVerificationFailed = errors.New("x402: payment verification failed")
	// ErrNetworkNotSupported indicates that the network is not supported.
	ErrNetworkNotSupported = errors.New("x402: network not supported")
)
