package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/dexfra-fun/x402-go/pkg/pricing"
	"github.com/dexfra-fun/x402-go/pkg/x402"
	chix402 "github.com/dexfra-fun/x402-go/pkg/adapters/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// Standard Chi middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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

	// Configure x402 middleware with path-based pricing
	config := &x402.Config{
		RecipientAddress: recipientAddress,
		Network:          network,
		FacilitatorURL:   facilitatorURL,
		PricingStrategy: pricing.NewPathBased(map[string]float64{
			"/api/data":    0.001, // 0.001 USDC
			"/api/premium": 0.01,  // 0.01 USDC
			"/api/action":  0.005, // 0.005 USDC
		}, 0.001), // default: 0.001 USDC
	}

	// Free endpoint - no payment required
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	// Protected routes group
	r.Route("/api", func(r chi.Router) {
		// Apply x402 middleware to all /api routes
		r.Use(chix402.NewMiddleware(config))

		r.Get("/data", func(w http.ResponseWriter, r *http.Request) {
			// Get payment info if needed
			if paymentInfo, ok := chix402.GetPaymentInfo(r.Context()); ok {
				log.Printf("Payment received: %.6f %s from %s",
					paymentInfo.Amount, paymentInfo.Currency, paymentInfo.Recipient)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "This is protected data",
				"data": map[string]interface{}{
					"temperature": 25.5,
					"humidity":    60,
					"timestamp":   "2025-01-01T00:00:00Z",
				},
			})
		})

		r.Get("/premium", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "This is premium content",
				"content": "Secret information only available after payment",
			})
		})

		r.Post("/action", func(w http.ResponseWriter, r *http.Request) {
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Action performed successfully",
				"result":  body,
			})
		})
	})

	log.Printf("Starting server on :8080")
	log.Printf("Network: %s", network)
	log.Printf("Recipient: %s", recipientAddress)
	log.Printf("Facilitator: %s", facilitatorURL)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
