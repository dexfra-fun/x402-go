// Package fiber provides x402 payment middleware adapter for Fiber framework.
package fiber

import (
	"github.com/dexfra-fun/x402-go/internal/common"
	"github.com/dexfra-fun/x402-go/pkg/x402"
	"github.com/gofiber/fiber/v2"
)

const paymentInfoKey = "x402_payment_info"

// NewMiddleware creates a new Fiber middleware for x402 payment handling.
func NewMiddleware(config *x402.Config) fiber.Handler {
	// Create common handler
	handler, err := common.NewHandler(config)
	if err != nil {
		config.Logger.Errorf("[x402-fiber] Failed to create middleware: %v", err)
		return func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Payment middleware configuration error",
			})
		}
	}

	return func(c *fiber.Ctx) error {
		// Extract resource from Fiber context
		resource := x402.Resource{
			Path:   c.Path(),
			Method: c.Method(),
			Params: make(map[string]string),
		}

		// Extract query parameters
		c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
			resource.Params[string(key)] = string(value)
		})

		// Get payment header
		paymentHeader := string(c.Request().Header.Peek("X-PAYMENT"))

		// Process payment
		result := handler.ProcessPaymentWithHeader(c.Context(), resource, paymentHeader)

		// Handle errors
		if result.Error != nil {
			return c.Status(result.StatusCode).JSON(fiber.Map{
				"error": result.ErrorMessage,
			})
		}

		// Handle payment required
		if result.RequirementNeeded {
			// Return proper x402 format response
			response := map[string]any{
				"x402Version": 1,
				"error":       "Payment required for this resource",
				"accepts":     []any{result.Requirement},
			}

			c.Set("Content-Type", "application/json")
			return c.Status(fiber.StatusPaymentRequired).JSON(response)
		}

		// Store payment info in context
		if result.PaymentInfo != nil {
			c.Locals(paymentInfoKey, result.PaymentInfo)
		}

		// Add X-PAYMENT-RESPONSE header if settlement was successful
		if result.Settlement != nil {
			encoded, err := x402.EncodeSettlement(*result.Settlement)
			if err != nil {
				config.Logger.Errorf("[x402-fiber] Failed to encode settlement: %v", err)
			} else {
				c.Set("X-Payment-Response", encoded)
			}
		}

		// Payment verified (or free endpoint) - proceed with request
		return c.Next()
	}
}

// GetPaymentInfo retrieves payment information from the Fiber context.
func GetPaymentInfo(c *fiber.Ctx) (*x402.PaymentInfo, bool) {
	if info := c.Locals(paymentInfoKey); info != nil {
		if paymentInfo, ok := info.(*x402.PaymentInfo); ok {
			return paymentInfo, true
		}
	}
	return nil, false
}
