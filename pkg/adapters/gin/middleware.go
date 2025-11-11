// Package gin provides x402 payment middleware adapter for Gin framework.
package gin

import (
	"net/http"

	"github.com/dexfra-fun/x402-go/internal/common"
	localx402 "github.com/dexfra-fun/x402-go/pkg/x402"
	"github.com/gin-gonic/gin"
)

const paymentInfoKey = "x402_payment_info"

// NewMiddleware creates a new Gin middleware for x402 payment handling.
func NewMiddleware(config *localx402.Config) gin.HandlerFunc {
	// Create common handler
	handler, err := common.NewHandler(config)
	if err != nil {
		config.Logger.Errorf("[x402-gin] Failed to create middleware: %v", err)
		// Return a middleware that always returns error
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Payment middleware configuration error",
			})
			c.Abort()
		}
	}

	return func(c *gin.Context) {
		// Extract resource from request
		resource := common.ExtractResource(c.Request)

		// Process payment
		result := handler.ProcessPayment(c.Request.Context(), resource, c.Request)

		// Handle errors
		if result.Error != nil {
			c.JSON(result.StatusCode, gin.H{
				"error": result.ErrorMessage,
			})
			c.Abort()
			return
		}

		// Handle payment required
		if result.RequirementNeeded {
			if writeErr := localx402.WritePaymentRequired(c.Writer, *result.Requirement); writeErr != nil {
				config.Logger.Errorf("[x402-gin] Failed to write payment required: %v", writeErr)
			}
			c.Abort()
			return
		}

		// Store payment info in context
		if result.PaymentInfo != nil {
			c.Set(paymentInfoKey, result.PaymentInfo)
		}

		// Payment verified (or free endpoint) - proceed with request
		c.Next()
	}
}

// GetPaymentInfo retrieves payment information from the Gin context.
func GetPaymentInfo(c *gin.Context) (*localx402.PaymentInfo, bool) {
	if info, exists := c.Get(paymentInfoKey); exists {
		if paymentInfo, ok := info.(*localx402.PaymentInfo); ok {
			return paymentInfo, true
		}
	}
	return nil, false
}
