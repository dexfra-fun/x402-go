package facilitator

import (
	"context"

	x402 "github.com/dexfra-fun/x402-go"
)

// Interface defines the standard facilitator contract for payment verification and settlement.
type Interface interface {
	Verify(ctx context.Context, payment x402.PaymentPayload, requirement x402.PaymentRequirement) (*VerifyResponse, error)
	Settle(ctx context.Context, payment x402.PaymentPayload, requirement x402.PaymentRequirement) (*x402.SettlementResponse, error)
	Supported(ctx context.Context) (*SupportedResponse, error)
}

// VerifyResponse contains the payment verification result from the facilitator.
type VerifyResponse struct {
	IsValid       bool   `json:"isValid"`
	InvalidReason string `json:"invalidReason,omitempty"`
	Payer         string `json:"payer"`
}

// SupportedKind describes a supported payment type with its configuration.
type SupportedKind struct {
	X402Version int            `json:"x402Version"`
	Scheme      string         `json:"scheme"`
	Network     string         `json:"network"`
	Extra       map[string]any `json:"extra,omitempty"`
}

// SupportedResponse lists all payment types supported by the facilitator.
type SupportedResponse struct {
	Kinds []SupportedKind `json:"kinds"`
}
