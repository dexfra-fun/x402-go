// Package main demonstrates how to use x402 payment middleware with Fiber framework.
package main

import (
	"errors"
	"log"
	"os"

	"github.com/bytedance/sonic"
	fiberx402 "github.com/dexfra-fun/x402-go/pkg/adapters/fiber"
	"github.com/dexfra-fun/x402-go/pkg/pricing"
	"github.com/dexfra-fun/x402-go/pkg/x402"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/shopspring/decimal"
)

const (
	serverPort        = ":8080"
	sampleTemperature = 25.5
	sampleHumidity    = 60
)

func getConfig() (*x402.Config, error) {
	recipientAddress := os.Getenv("X402_RECIPIENT_ADDRESS")
	if recipientAddress == "" {
		return nil, errors.New("X402_RECIPIENT_ADDRESS environment variable is required")
	}

	network := os.Getenv("X402_NETWORK")
	if network == "" {
		network = "solana-devnet" // default
	}

	facilitatorURL := os.Getenv("X402_FACILITATOR_URL")
	if facilitatorURL == "" {
		facilitatorURL = "https://facilitator.payai.network" // default
	}

	return &x402.Config{
		RecipientAddress: recipientAddress,
		Network:          network,
		FacilitatorURL:   facilitatorURL,
		PricingStrategy: pricing.NewPathBased(map[string]decimal.Decimal{
			"/api/data":    decimal.RequireFromString("0.001"), // 0.001 USDC
			"/api/premium": decimal.RequireFromString("0.01"),  // 0.01 USDC
			"/api/action":  decimal.RequireFromString("0.005"), // 0.005 USDC
		}, decimal.RequireFromString("0.001")), // default: 0.001 USDC
	}, nil
}

func setupRoutes(app *fiber.App, config *x402.Config) {
	// Free endpoint - no payment required
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Protected routes group
	api := app.Group("/api")
	api.Use(fiberx402.NewMiddleware(config))

	api.Get("/data", func(c *fiber.Ctx) error {
		// Get payment info if needed
		if paymentInfo, ok := fiberx402.GetPaymentInfo(c); ok {
			log.Printf("Payment received: %s %s from %s",
				paymentInfo.Amount.String(), paymentInfo.Currency, paymentInfo.Recipient)
		}

		return c.JSON(fiber.Map{
			"message": "This is protected data",
			"data": fiber.Map{
				"temperature": sampleTemperature,
				"humidity":    sampleHumidity,
				"timestamp":   "2025-01-01T00:00:00Z",
			},
		})
	})

	api.Get("/premium", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "This is premium content",
			"content": "Secret information only available after payment",
		})
	})

	api.Post("/action", func(c *fiber.Ctx) error {
		var body map[string]any
		if err := sonic.ConfigDefault.Unmarshal(c.Body(), &body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Action performed successfully",
			"result":  body,
		})
	})
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "x402-go Fiber Example",
	})

	// Standard Fiber middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Get configuration and setup routes
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	setupRoutes(app, config)

	log.Printf("Starting server on %s", serverPort)
	log.Printf("Network: %s", config.Network)
	log.Printf("Recipient: %s", config.RecipientAddress)
	log.Printf("Facilitator: %s", config.FacilitatorURL)

	if err := app.Listen(serverPort); err != nil {
		log.Fatal(err)
	}
}
