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
		paymentHeader := string(c.Request().Header.Peek("X-402-Payment"))

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
			// Set payment requirement headers
			c.Set("X-402-Version", "1")
			c.Set("X-402-Network", config.Network)
			c.Set("X-402-Recipient", config.RecipientAddress)
			if result.PaymentInfo != nil {
				c.Set("X-402-Amount", result.PaymentInfo.Amount.String())
				c.Set("X-402-Currency", result.PaymentInfo.Currency)
			}

			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
				"error":   "Payment required",
				"network": config.Network,
			})
		}

		// Store payment info in context
		if result.PaymentInfo != nil {
			c.Locals(paymentInfoKey, result.PaymentInfo)
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
