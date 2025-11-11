package main

import (
	"log"
	"os"

	"github.com/dexfra-fun/x402-go/pkg/pricing"
	"github.com/dexfra-fun/x402-go/pkg/x402"
	ginx402 "github.com/dexfra-fun/x402-go/pkg/adapters/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Get configuration from environment variables
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

	// Configure x402 middleware with fixed pricing
	config := &x402.Config{
		RecipientAddress: recipientAddress,
		Network:          network,
		FacilitatorURL:   facilitatorURL,
		PricingStrategy:  pricing.NewFixed(0.001), // 0.001 USDC per call
	}

	// Free endpoint - no payment required
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
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
				log.Printf("Payment received: %.6f %s from %s",
					paymentInfo.Amount, paymentInfo.Currency, paymentInfo.Recipient)
			}

			c.JSON(200, gin.H{
				"message": "This is protected data",
				"data": map[string]interface{}{
					"temperature": 25.5,
					"humidity":    60,
					"timestamp":   "2025-01-01T00:00:00Z",
				},
			})
		})

		protected.GET("/premium", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "This is premium content",
				"content": "Secret information only available after payment",
			})
		})

		protected.POST("/action", func(c *gin.Context) {
			var body map[string]interface{}
			if err := c.BindJSON(&body); err != nil {
				c.JSON(400, gin.H{"error": "Invalid request body"})
				return
			}

			c.JSON(200, gin.H{
				"message": "Action performed successfully",
				"result":  body,
			})
		})
	}

	log.Printf("Starting server on :8080")
	log.Printf("Network: %s", network)
	log.Printf("Recipient: %s", recipientAddress)
	log.Printf("Facilitator: %s", facilitatorURL)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
