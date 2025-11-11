package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	httpx402 "github.com/dexfra-fun/x402-go/pkg/adapters/http"
	"github.com/dexfra-fun/x402-go/pkg/pricing"
	"github.com/dexfra-fun/x402-go/pkg/x402"
	"github.com/shopspring/decimal"
)

func main() {
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

	// Configure x402 middleware with method-based pricing
	config := &x402.Config{
		RecipientAddress: recipientAddress,
		Network:          network,
		FacilitatorURL:   facilitatorURL,
		PricingStrategy: pricing.NewMethodBased(map[string]decimal.Decimal{
			"GET":    decimal.RequireFromString("0.001"), // 0.001 USDC for reads
			"POST":   decimal.RequireFromString("0.005"), // 0.005 USDC for writes
			"PUT":    decimal.RequireFromString("0.005"), // 0.005 USDC for updates
			"DELETE": decimal.RequireFromString("0.01"),  // 0.01 USDC for deletes
		}, decimal.RequireFromString("0.001")), // default: 0.001 USDC
	}

	mux := http.NewServeMux()

	// Free endpoint - no payment required
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	// Protected endpoints
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		// Get payment info if needed
		if paymentInfo, ok := httpx402.GetPaymentInfo(r.Context()); ok {
			log.Printf("Payment received: %s %s from %s",
				paymentInfo.Amount.String(), paymentInfo.Currency, paymentInfo.Recipient)
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

	mux.HandleFunc("/api/premium", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "This is premium content",
			"content": "Secret information only available after payment",
		})
	})

	mux.HandleFunc("/api/action", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

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

	// Wrap mux with x402 middleware
	handler := httpx402.NewMiddleware(config)(mux)

	log.Printf("Starting server on :8080")
	log.Printf("Network: %s", network)
	log.Printf("Recipient: %s", recipientAddress)
	log.Printf("Facilitator: %s", facilitatorURL)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
