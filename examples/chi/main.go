// Package main demonstrates how to use x402 payment middleware with Chi router.
package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bytedance/sonic"
	chix402 "github.com/dexfra-fun/x402-go/pkg/adapters/chi"
	"github.com/dexfra-fun/x402-go/pkg/pricing"
	"github.com/dexfra-fun/x402-go/pkg/x402"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func setupRoutes(r chi.Router, config *x402.Config) {
	// Free endpoint - no payment required
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := sonic.ConfigDefault.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	// Protected routes group
	r.Route("/api", func(r chi.Router) {
		// Apply x402 middleware to all /api routes
		r.Use(chix402.NewMiddleware(config))

		r.Get("/data", func(w http.ResponseWriter, r *http.Request) {
			// Get payment info if needed
			if paymentInfo, ok := chix402.GetPaymentInfo(r.Context()); ok {
				log.Printf("Payment received: %s %s from %s",
					paymentInfo.Amount.String(), paymentInfo.Currency, paymentInfo.Recipient)
			}

			w.Header().Set("Content-Type", "application/json")
			if err := sonic.ConfigDefault.NewEncoder(w).Encode(map[string]any{
				"message": "This is protected data",
				"data": map[string]any{
					"temperature": sampleTemperature,
					"humidity":    sampleHumidity,
					"timestamp":   "2025-01-01T00:00:00Z",
				},
			}); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		})

		r.Get("/premium", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := sonic.ConfigDefault.NewEncoder(w).Encode(map[string]any{
				"message": "This is premium content",
				"content": "Secret information only available after payment",
			}); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		})

		r.Post("/action", func(w http.ResponseWriter, r *http.Request) {
			var body map[string]any
			if err := sonic.ConfigDefault.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if err := sonic.ConfigDefault.NewEncoder(w).Encode(map[string]any{
				"message": "Action performed successfully",
				"result":  body,
			}); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		})
	})
}

func main() {
	r := chi.NewRouter()

	// Standard Chi middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Get configuration and setup routes
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	setupRoutes(r, config)

	log.Printf("Starting server on %s", serverPort)
	log.Printf("Network: %s", config.Network)
	log.Printf("Recipient: %s", config.RecipientAddress)
	log.Printf("Facilitator: %s", config.FacilitatorURL)

	// Create server with timeouts for security
	server := &http.Server{
		Addr:              serverPort,
		Handler:           r,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
