package x402

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
	x402 "github.com/dexfra-fun/x402-go"
)

const (
	// HeaderPayment is the HTTP header name for x402 payment information.
	// Uses canonical form for HTTP headers.
	HeaderPayment = "X-Payment"
	// HeaderPaymentResponse is the HTTP header name for x402 settlement response.
	// Uses canonical form for HTTP headers.
	HeaderPaymentResponse = "X-Payment-Response"
)

// EncodePaymentRequirement encodes a payment requirement as a base64 JSON string.
func EncodePaymentRequirement(req x402.PaymentRequirement) (string, error) {
	jsonBytes, err := sonic.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal payment requirement: %w", err)
	}
	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

// DecodePaymentRequirement decodes a base64 JSON payment requirement string.
func DecodePaymentRequirement(encoded string) (*x402.PaymentRequirement, error) {
	jsonBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("decode base64: %w", err)
	}

	var req x402.PaymentRequirement
	if err := sonic.Unmarshal(jsonBytes, &req); err != nil {
		return nil, fmt.Errorf("unmarshal payment requirement: %w", err)
	}

	return &req, nil
}

// EncodePaymentPayload encodes a payment payload as a base64 JSON string.
func EncodePaymentPayload(payload x402.PaymentPayload) (string, error) {
	jsonBytes, err := sonic.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payment payload: %w", err)
	}
	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

// DecodePaymentPayload decodes a base64 JSON payment payload string.
func DecodePaymentPayload(encoded string) (*x402.PaymentPayload, error) {
	jsonBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("decode base64: %w", err)
	}

	var payload x402.PaymentPayload
	if err := sonic.Unmarshal(jsonBytes, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal payment payload: %w", err)
	}

	return &payload, nil
}

// ParsePaymentHeader extracts and decodes the payment payload from an HTTP request.
func ParsePaymentHeader(r *http.Request) (*x402.PaymentPayload, error) {
	header := r.Header.Get(HeaderPayment)
	if header == "" {
		return nil, fmt.Errorf("missing %s header", HeaderPayment)
	}

	return DecodePaymentPayload(header)
}

// SetPaymentResponseHeader sets the X-PAYMENT-RESPONSE header with settlement information.
func SetPaymentResponseHeader(w http.ResponseWriter, settlement x402.SettlementResponse) error {
	encoded, err := EncodeSettlement(settlement)
	if err != nil {
		return err
	}
	w.Header().Set(HeaderPaymentResponse, encoded)
	return nil
}

// SetPaymentRequiredHeader sets payment requirement in the response body (legacy function for 402 responses).
// Note: This should use WritePaymentRequired instead for proper x402 format.
func SetPaymentRequiredHeader(w http.ResponseWriter, req x402.PaymentRequirement) error {
	// For backward compatibility, but WritePaymentRequired is preferred
	return WritePaymentRequired(w, req)
}

// WritePaymentRequired writes a 402 Payment Required response with proper x402 format.
func WritePaymentRequired(w http.ResponseWriter, req x402.PaymentRequirement) error {
	// Create proper x402 response body according to specification
	response := x402.PaymentRequirementsResponse{
		X402Version: 1,
		Error:       "Payment required for this resource",
		Accepts:     []x402.PaymentRequirement{req},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusPaymentRequired)
	return sonic.ConfigDefault.NewEncoder(w).Encode(response)
}

// EncodeSettlement encodes a settlement response as a base64 JSON string.
func EncodeSettlement(settlement x402.SettlementResponse) (string, error) {
	jsonBytes, err := sonic.Marshal(settlement)
	if err != nil {
		return "", fmt.Errorf("marshal settlement response: %w", err)
	}
	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

// DecodeSettlement decodes a base64 JSON settlement response string.
func DecodeSettlement(encoded string) (*x402.SettlementResponse, error) {
	jsonBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("decode base64: %w", err)
	}

	var settlement x402.SettlementResponse
	if err := sonic.Unmarshal(jsonBytes, &settlement); err != nil {
		return nil, fmt.Errorf("unmarshal settlement response: %w", err)
	}

	return &settlement, nil
}

// BasicPaymentCheck performs basic validation that a payment matches a requirement.
// Note: Full payment verification must be done by the facilitator.
func BasicPaymentCheck(payload x402.PaymentPayload, req x402.PaymentRequirement) bool {
	// Check scheme
	if !strings.EqualFold(payload.Scheme, req.Scheme) {
		return false
	}

	// Check network
	if !strings.EqualFold(payload.Network, req.Network) {
		return false
	}

	return true
}
