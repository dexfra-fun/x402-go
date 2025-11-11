// Package http provides x402 payment middleware adapter for standard HTTP handlers.
package http

import (
	"context"
	"net/http"

	"github.com/dexfra-fun/x402-go/pkg/x402"
	mark3labs "github.com/mark3labs/x402-go"
	x402http "github.com/mark3labs/x402-go/http"
)

type contextKey string

const paymentInfoKey contextKey = "x402_payment_info"

// NewMiddleware creates a new standard HTTP middleware for x402 payment handling.
func NewMiddleware(config *x402.Config) func(http.Handler) http.Handler {
	// Create x402 middleware
	middleware, err := x402.New(config)
	if err != nil {
		config.Logger.Errorf("[x402-http] Failed to create middleware: %v", err)
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
			// Create resource from request
			resource := x402.Resource{
				Path:   r.URL.Path,
				Method: r.Method,
				Params: make(map[string]string),
			}

			// Extract query parameters
			for key, values := range r.URL.Query() {
				if len(values) > 0 {
					resource.Params[key] = values[0]
				}
			}

			// Process payment requirement
			requirement, paymentInfo, err := middleware.ProcessRequest(r.Context(), resource)
			if err != nil {
				config.Logger.Errorf("[x402-http] Failed to process payment: %v", err)
				http.Error(w, "Payment processing error", http.StatusInternalServerError)
				return
			}

			// Free endpoint - no payment required
			if requirement == nil {
				next.ServeHTTP(w, r)
				return
			}

			// Payment required - configure x402 HTTP middleware
			x402Config := &x402http.Config{
				FacilitatorURL: config.FacilitatorURL,
				PaymentRequirements: []mark3labs.PaymentRequirement{
					*requirement,
				},
			}

			// Apply mark3labs x402 HTTP middleware
			x402Handler := x402http.NewX402Middleware(x402Config)(next)

			// Store payment info in request context for later use
			if paymentInfo != nil {
				ctx := r.Context()
				ctx = context.WithValue(ctx, paymentInfoKey, paymentInfo)
				r = r.WithContext(ctx)
			}

			x402Handler.ServeHTTP(w, r)
		})
	}
}

// GetPaymentInfo retrieves payment information from the request context.
func GetPaymentInfo(ctx context.Context) (*x402.PaymentInfo, bool) {
	info, ok := ctx.Value(paymentInfoKey).(*x402.PaymentInfo)
	return info, ok
}
