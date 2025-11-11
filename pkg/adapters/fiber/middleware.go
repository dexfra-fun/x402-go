// Package fiber provides x402 payment middleware adapter for Fiber framework.
package fiber

import (
	"github.com/dexfra-fun/x402-go/pkg/x402"
	"github.com/gofiber/fiber/v2"
)

// createResource creates an x402 resource from a Fiber context.
func createResource(c *fiber.Ctx) x402.Resource {
	resource := x402.Resource{
		Path:   c.Path(),
		Method: c.Method(),
		Params: make(map[string]string),
	}

	// Extract query parameters
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		resource.Params[string(key)] = string(value)
	})

	return resource
}

// sendPaymentRequired sends a 402 Payment Required response.
func sendPaymentRequired(c *fiber.Ctx, config *x402.Config, paymentInfo *x402.PaymentInfo) error {
	c.Set("X-402-Version", "1")
	c.Set("X-402-Methods", "solana")
	c.Set("X-402-Network", config.Network)
	c.Set("X-402-Recipient", config.RecipientAddress)
	c.Set("X-402-Amount", paymentInfo.Amount.String())
	c.Set("X-402-Currency", paymentInfo.Currency)

	return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
		"error": "Payment required",
		"payment": fiber.Map{
			"version":   1,
			"methods":   []string{"solana"},
			"network":   config.Network,
			"recipient": config.RecipientAddress,
			"amount":    paymentInfo.Amount.String(),
			"currency":  paymentInfo.Currency,
		},
	})
}

// NewMiddleware creates a new Fiber middleware for x402 payment handling.
func NewMiddleware(config *x402.Config) fiber.Handler {
	middleware, err := x402.New(config)
	if err != nil {
		config.Logger.Errorf("[x402-fiber] Failed to create middleware: %v", err)
		return func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Payment middleware configuration error",
			})
		}
	}

	return func(c *fiber.Ctx) error {
		resource := createResource(c)
		requirement, paymentInfo, err := middleware.ProcessRequest(c.Context(), resource)
		if err != nil {
			config.Logger.Errorf("[x402-fiber] Failed to process payment: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Payment processing error",
			})
		}

		if requirement == nil {
			return c.Next()
		}

		paymentHeader := string(c.Request().Header.Peek("X-Payment"))
		if paymentHeader == "" {
			return sendPaymentRequired(c, config, paymentInfo)
		}

		if paymentInfo != nil {
			c.Locals("x402_payment_info", paymentInfo)
		}

		return c.Next()
	}
}

// GetPaymentInfo retrieves payment information from the Fiber context.
func GetPaymentInfo(c *fiber.Ctx) (*x402.PaymentInfo, bool) {
	if info := c.Locals("x402_payment_info"); info != nil {
		if paymentInfo, ok := info.(*x402.PaymentInfo); ok {
			return paymentInfo, true
		}
	}
	return nil, false
}
