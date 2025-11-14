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
	HeaderPayment = "X-402-Payment"
	// HeaderPaymentRequired is the HTTP header name for x402 payment requirements.
	HeaderPaymentRequired = "X-402-Payment-Required"
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

// SetPaymentRequiredHeader sets the X-402-Payment-Required header on the response.
func SetPaymentRequiredHeader(w http.ResponseWriter, req x402.PaymentRequirement) error {
	encoded, err := EncodePaymentRequirement(req)
	if err != nil {
		return err
	}
	w.Header().Set(HeaderPaymentRequired, encoded)
	return nil
}

// WritePaymentRequired writes a 402 Payment Required response with proper x402 format.
func WritePaymentRequired(w http.ResponseWriter, req x402.PaymentRequirement) error {
	if err := SetPaymentRequiredHeader(w, req); err != nil {
		return err
	}

	// Create proper x402 response body according to specification
	response := x402.PaymentRequirementsResponse{
		X402Version: 1,
		Error:       "X-PAYMENT header is required",
		Accepts:     []x402.PaymentRequirement{req},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusPaymentRequired)
	return sonic.ConfigDefault.NewEncoder(w).Encode(response)
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
