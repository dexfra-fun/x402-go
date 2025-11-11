// Package chi provides x402 payment middleware adapter for Chi router.
package chi

import (
	"context"
	"net/http"

	"github.com/dexfra-fun/x402-go/internal/common"
	localx402 "github.com/dexfra-fun/x402-go/pkg/x402"
)

type contextKey string

const paymentInfoKey contextKey = "x402_payment_info"

// NewMiddleware creates a new Chi middleware for x402 payment handling.
func NewMiddleware(config *localx402.Config) func(http.Handler) http.Handler {
	// Create common handler
	handler, err := common.NewHandler(config)
	if err != nil {
		config.Logger.Errorf("[x402-chi] Failed to create middleware: %v", err)
		// Return a middleware that always returns error
		return func(_ http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(
					w,
					"Payment middleware configuration error",
					http.StatusInternalServerError,
				)
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract resource from request
			resource := common.ExtractResource(r)

			// Process payment
			result := handler.ProcessPayment(r.Context(), resource, r)

			// Handle errors
			if result.Error != nil {
				http.Error(w, result.ErrorMessage, result.StatusCode)
				return
			}

			// Handle payment required
			if result.RequirementNeeded {
				if writeErr := localx402.WritePaymentRequired(w, *result.Requirement); writeErr != nil {
					config.Logger.Errorf("[x402-chi] Failed to write payment required: %v", writeErr)
				}
				return
			}

			// Store payment info in request context
			if result.PaymentInfo != nil {
				ctx := r.Context()
				ctx = context.WithValue(ctx, paymentInfoKey, result.PaymentInfo)
				r = r.WithContext(ctx)
			}

			// Payment verified (or free endpoint) - proceed with request
			next.ServeHTTP(w, r)
		})
	}
}

// GetPaymentInfo retrieves payment information from the request context.
func GetPaymentInfo(ctx context.Context) (*localx402.PaymentInfo, bool) {
	info, ok := ctx.Value(paymentInfoKey).(*localx402.PaymentInfo)
	return info, ok
}
