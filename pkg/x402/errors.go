package x402

import "errors"

var (
	// Configuration errors.
	ErrMissingRecipient   = errors.New("x402: recipient address is required")
	ErrMissingNetwork     = errors.New("x402: network is required")
	ErrMissingFacilitator = errors.New("x402: facilitator URL is required")
	ErrMissingPricing     = errors.New("x402: pricing strategy is required")

	// Runtime errors.
	ErrInvalidFeePayer           = errors.New("x402: invalid fee payer address")
	ErrFacilitatorUnavailable    = errors.New("x402: facilitator service unavailable")
	ErrFeePayerNotFound          = errors.New("x402: fee payer not found for network")
	ErrPaymentVerificationFailed = errors.New("x402: payment verification failed")
	ErrNetworkNotSupported       = errors.New("x402: network not supported")
)
