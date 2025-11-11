// Package main demonstrates how to use x402 payment middleware with standard HTTP handlers.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	httpx402 "github.com/dexfra-fun/x402-go/pkg/adapters/http"
	"github.com/dexfra-fun/x402-go/pkg/pricing"
	"github.com/dexfra-fun/x402-go/pkg/x402"
	"github.com/shopspring/decimal"
)

const (
	serverPort        = ":8080"
	readTimeout       = 15 * time.Second
	readHeaderTimeout = 10 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	sampleTemperature = 25.5
	sampleHumidity    = 60
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
		PricingStrategy: pricing.NewMethodBased(map[string]decimal.Decimal{
			"GET":    decimal.RequireFromString("0.001"), // 0.001 USDC for reads
			"POST":   decimal.RequireFromString("0.005"), // 0.005 USDC for writes
			"PUT":    decimal.RequireFromString("0.005"), // 0.005 USDC for updates
			"DELETE": decimal.RequireFromString("0.01"),  // 0.01 USDC for deletes
		}, decimal.RequireFromString("0.001")), // default: 0.001 USDC
	}
}

func setupHandlers(mux *http.ServeMux) {
	// Free endpoint - no payment required
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	// Protected endpoints
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		// Get payment info if needed
		if paymentInfo, ok := httpx402.GetPaymentInfo(r.Context()); ok {
			log.Printf("Payment received: %s %s from %s",
				paymentInfo.Amount.String(), paymentInfo.Currency, paymentInfo.Recipient)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "This is protected data",
			"data": map[string]interface{}{
				"temperature": sampleTemperature,
				"humidity":    sampleHumidity,
				"timestamp":   "2025-01-01T00:00:00Z",
			},
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/api/premium", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "This is premium content",
			"content": "Secret information only available after payment",
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
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
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Action performed successfully",
			"result":  body,
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})
}

func main() {
	config := getConfig()

	mux := http.NewServeMux()
	setupHandlers(mux)

	// Wrap mux with x402 middleware
	handler := httpx402.NewMiddleware(config)(mux)

	log.Printf("Starting server on %s", serverPort)
	log.Printf("Network: %s", config.Network)
	log.Printf("Recipient: %s", config.RecipientAddress)
	log.Printf("Facilitator: %s", config.FacilitatorURL)

	// Create server with timeouts for security
	server := &http.Server{
		Addr:              serverPort,
		Handler:           handler,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
