// Package main demonstrates how to use x402 payment middleware with Gin framework.
package main

import (
	"log"
	"os"

	ginx402 "github.com/dexfra-fun/x402-go/pkg/adapters/gin"
	"github.com/dexfra-fun/x402-go/pkg/pricing"
	"github.com/dexfra-fun/x402-go/pkg/x402"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

const (
	httpStatusOK         = 200
	httpStatusBadRequest = 400
	sampleTemperature    = 25.5
	sampleHumidity       = 60
)

func getConfig() *x402.Config {
	recipientAddress := os.Getenv("X402_RECIPIENT_ADDRESS")
	if recipientAddress == "" {
		log.Fatal("X402_RECIPIENT_ADDRESS environment variable is required")
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
		PricingStrategy:  pricing.NewFixed(decimal.RequireFromString("0.001")),
	}
}

func setupRoutes(r *gin.Engine, config *x402.Config) {
	// Free endpoint - no payment required
	r.GET("/health", func(c *gin.Context) {
		c.JSON(httpStatusOK, gin.H{
			"status": "ok",
		})
	})

	// Protected endpoints group
	protected := r.Group("/api")
	protected.Use(ginx402.NewMiddleware(config))
	{
		protected.GET("/data", func(c *gin.Context) {
			// Get payment info if needed
			if paymentInfo, ok := ginx402.GetPaymentInfo(c); ok {
				log.Printf("Payment received: %s %s from %s",
					paymentInfo.Amount.String(), paymentInfo.Currency, paymentInfo.Recipient)
			}

			c.JSON(httpStatusOK, gin.H{
				"message": "This is protected data",
				"data": map[string]interface{}{
					"temperature": sampleTemperature,
					"humidity":    sampleHumidity,
					"timestamp":   "2025-01-01T00:00:00Z",
				},
			})
		})

		protected.GET("/premium", func(c *gin.Context) {
			c.JSON(httpStatusOK, gin.H{
				"message": "This is premium content",
				"content": "Secret information only available after payment",
			})
		})

		protected.POST("/action", func(c *gin.Context) {
			var body map[string]interface{}
			if err := c.BindJSON(&body); err != nil {
				c.JSON(httpStatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			c.JSON(httpStatusOK, gin.H{
				"message": "Action performed successfully",
				"result":  body,
			})
		})
	}
}

func main() {
	r := gin.Default()

	// Get configuration and setup routes
	config := getConfig()
	setupRoutes(r, config)

	log.Printf("Starting server on :8080")
	log.Printf("Network: %s", config.Network)
	log.Printf("Recipient: %s", config.RecipientAddress)
	log.Printf("Facilitator: %s", config.FacilitatorURL)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
