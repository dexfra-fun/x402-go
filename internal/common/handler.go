// Package common provides shared functionality for x402 middleware adapters.
package common

import (
	"context"
	"net/http"

	x402 "github.com/dexfra-fun/x402-go"
	localx402 "github.com/dexfra-fun/x402-go/pkg/x402"
)

// PaymentResult represents the outcome of payment processing.
type PaymentResult struct {
	// RequirementNeeded indicates if a 402 Payment Required response should be sent
	RequirementNeeded bool
	// Requirement contains the payment requirement (if needed)
	Requirement *x402.PaymentRequirement
	// Error indicates if an error occurred
	Error error
	// ErrorMessage is the user-facing error message
	ErrorMessage string
	// StatusCode is the HTTP status code to return on error
	StatusCode int
	// PaymentInfo contains metadata about the payment
	PaymentInfo *localx402.PaymentInfo
	// Payer is the address of the payer (after successful verification)
	Payer string
}

// Handler encapsulates common payment processing logic.
type Handler struct {
	middleware *localx402.Middleware
	config     *localx402.Config
}

// NewHandler creates a new common payment handler.
func NewHandler(config *localx402.Config) (*Handler, error) {
	middleware, err := localx402.New(config)
	if err != nil {
		return nil, err
	}

	return &Handler{
		middleware: middleware,
		config:     config,
	}, nil
}

// ProcessPayment performs the complete payment processing flow.
// Returns PaymentResult indicating what action should be taken.
func (h *Handler) ProcessPayment(
	ctx context.Context,
	resource localx402.Resource,
	r *http.Request,
) PaymentResult {
	// Step 1: Get payment requirement
	requirement, paymentInfo, err := h.middleware.ProcessRequest(ctx, resource)
	if err != nil {
		h.config.Logger.Errorf("[x402-common] Failed to process payment: %v", err)
		return PaymentResult{
			Error:        err,
			ErrorMessage: "Payment processing error",
			StatusCode:   http.StatusInternalServerError,
		}
	}

	// Step 2: Check if payment is required
	if requirement == nil {
		// Free endpoint - no payment required
		return PaymentResult{
			RequirementNeeded: false,
			PaymentInfo:       paymentInfo,
		}
	}

	// Step 3: Check if payment was provided
	payment, err := localx402.ParsePaymentHeader(r)
	if err != nil {
		// No payment provided - return 402 with requirement
		h.config.Logger.Printf("[x402-common] No payment provided: %v", err)
		return PaymentResult{
			RequirementNeeded: true,
			Requirement:       requirement,
			PaymentInfo:       paymentInfo,
		}
	}

	return h.verifyAndSettle(ctx, payment, requirement, paymentInfo)
}

// ProcessPaymentWithHeader performs payment processing with payment header string.
// Useful for frameworks that don't use standard http.Request (e.g., Fiber with fasthttp).
func (h *Handler) ProcessPaymentWithHeader(
	ctx context.Context,
	resource localx402.Resource,
	paymentHeader string,
) PaymentResult {
	// Step 1: Get payment requirement
	requirement, paymentInfo, err := h.middleware.ProcessRequest(ctx, resource)
	if err != nil {
		h.config.Logger.Errorf("[x402-common] Failed to process payment: %v", err)
		return PaymentResult{
			Error:        err,
			ErrorMessage: "Payment processing error",
			StatusCode:   http.StatusInternalServerError,
		}
	}

	// Step 2: Check if payment is required
	if requirement == nil {
		// Free endpoint - no payment required
		return PaymentResult{
			RequirementNeeded: false,
			PaymentInfo:       paymentInfo,
		}
	}

	// Step 3: Check if payment was provided
	if paymentHeader == "" {
		return PaymentResult{
			RequirementNeeded: true,
			Requirement:       requirement,
			PaymentInfo:       paymentInfo,
		}
	}

	// Step 4: Decode payment
	payment, err := localx402.DecodePaymentPayload(paymentHeader)
	if err != nil {
		h.config.Logger.Printf("[x402-common] Failed to decode payment: %v", err)
		return PaymentResult{
			RequirementNeeded: true,
			Requirement:       requirement,
			PaymentInfo:       paymentInfo,
		}
	}

	return h.verifyAndSettle(ctx, payment, requirement, paymentInfo)
}

// verifyAndSettle performs payment verification and settlement.
func (h *Handler) verifyAndSettle(
	ctx context.Context,
	payment *x402.PaymentPayload,
	requirement *x402.PaymentRequirement,
	paymentInfo *localx402.PaymentInfo,
) PaymentResult {
	// Step 1: Basic validation
	if !localx402.BasicPaymentCheck(*payment, *requirement) {
		h.config.Logger.Errorf("[x402-common] Payment does not match requirement")
		return PaymentResult{
			Error:        localx402.ErrPaymentVerificationFailed,
			ErrorMessage: "Invalid payment",
			StatusCode:   http.StatusBadRequest,
		}
	}

	// Step 2: Verify payment with facilitator
	isValid, invalidReason, payer, err := h.middleware.GetFacilitator().Verify(
		ctx,
		*payment,
		*requirement,
	)
	if err != nil {
		h.config.Logger.Errorf("[x402-common] Failed to verify payment: %v", err)
		return PaymentResult{
			Error:        err,
			ErrorMessage: "Payment verification error",
			StatusCode:   http.StatusInternalServerError,
		}
	}

	if !isValid {
		h.config.Logger.Errorf("[x402-common] Payment verification failed: %s", invalidReason)
		return PaymentResult{
			Error:        localx402.ErrPaymentVerificationFailed,
			ErrorMessage: "Payment verification failed: " + invalidReason,
			StatusCode:   http.StatusPaymentRequired,
		}
	}

	h.config.Logger.Printf("[x402-common] Payment verified: payer=%s", payer)

	// Step 3: Settle payment asynchronously
	go h.settlePaymentAsync(*payment, *requirement)

	// Step 4: Return success
	return PaymentResult{
		RequirementNeeded: false,
		PaymentInfo:       paymentInfo,
		Payer:             payer,
	}
}

// settlePaymentAsync settles payment in the background.
func (h *Handler) settlePaymentAsync(
	payment x402.PaymentPayload,
	requirement x402.PaymentRequirement,
) {
	settlementCtx := context.Background()
	settlement, err := h.middleware.GetFacilitator().Settle(
		settlementCtx,
		payment,
		requirement,
	)
	if err != nil {
		h.config.Logger.Errorf("[x402-common] Failed to settle payment: %v", err)
		return
	}
	if !settlement.Success {
		h.config.Logger.Errorf("[x402-common] Settlement failed: %s", settlement.ErrorReason)
		return
	}
	h.config.Logger.Printf("[x402-common] Payment settled: tx=%s", settlement.Transaction)
}

// GetConfig returns the handler configuration.
func (h *Handler) GetConfig() *localx402.Config {
	return h.config
}

// GetMiddleware returns the underlying middleware instance.
func (h *Handler) GetMiddleware() *localx402.Middleware {
	return h.middleware
}
