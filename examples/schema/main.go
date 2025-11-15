// Package main demonstrates how to use schema support with x402 payment middleware.
// This example shows how to define input/output schemas for API endpoints according to
// the official x402 specification.
package main

import (
	"errors"
	"log"
	"os"

	x402 "github.com/dexfra-fun/x402-go"
	ginx402 "github.com/dexfra-fun/x402-go/pkg/adapters/gin"
	"github.com/dexfra-fun/x402-go/pkg/pricing"
	"github.com/dexfra-fun/x402-go/pkg/resource"
	"github.com/dexfra-fun/x402-go/pkg/schema"
	localx402 "github.com/dexfra-fun/x402-go/pkg/x402"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

const (
	httpStatusOK         = 200
	httpStatusBadRequest = 400
	defaultResultCount   = 2
	sampleTemperature    = 25.5
)

func createTwitterFollowersSchema() *x402.EndpointSchema {
	return &x402.EndpointSchema{
		Input: &x402.InputSchema{
			Type:   "http",
			Method: "GET",
			QueryParams: map[string]*x402.FieldDef{
				"userName": {
					Type:        "string",
					Required:    true,
					Description: "screen name of the user",
				},
				"cursor": {
					Type:        "string",
					Required:    false,
					Description: "The cursor to paginate through the results. First page is empty string.",
				},
				"pageSize": {
					Type:        "integer",
					Required:    false,
					Description: "The number of followings to return per page. Default is 200. min 20, max 200",
				},
			},
			HeaderFields: map[string]*x402.FieldDef{
				"aisa-payment-token": {
					Type:        "string",
					Required:    false,
					Description: "Optional merchant payment token for wallet-based payment mode",
				},
				"aisa-payment": {
					Type:        "string",
					Required:    false,
					Description: "Payment proof header (base64 encoded JSON) for signature-based payment mode",
				},
			},
		},
		Output: map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id":         map[string]any{"type": "string"},
					"username":   map[string]any{"type": "string"},
					"followedAt": map[string]any{"type": "string", "format": "date-time"},
				},
			},
		},
	}
}

func createDataSubmitSchema() *x402.EndpointSchema {
	return schema.NewEndpointSchema().
		WithInput(
			schema.NewInputSchema("POST").
				WithBodyType("json").
				WithBodyField("query", schema.NewFieldDef("string", true, "Search query")).
				WithBodyField("limit", schema.NewFieldDef("integer", false, "Result limit (default: 10)")).
				WithBodyField("filters", schema.NewObjectField(map[string]*x402.FieldDef{
					"category": schema.NewFieldDef("string", false, "Filter by category"),
					"status":   schema.NewEnumField([]string{"active", "inactive", "pending"}, false, "Filter by status"),
				}, false, "Optional filters")).
				WithHeaderField("X-API-Key", schema.NewFieldDef("string", false, "Optional API key for enhanced access")).
				Build(),
		).
		WithOutput(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"results": map[string]any{"type": "array"},
				"total":   map[string]any{"type": "integer"},
			},
		}).
		Build()
}

func createWeatherSchema() *x402.EndpointSchema {
	return schema.NewEndpointSchema().
		WithInput(
			schema.NewInputSchema("GET").
				WithQueryParam("city", schema.NewFieldDef("string", true, "City name")).
				WithQueryParam("units", schema.NewEnumField([]string{"metric", "imperial"}, false, "Temperature units")).
				Build(),
		).
		Build()
}

func getConfig() (*localx402.Config, error) {
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

	// Get base URL for resource URLs (optional)
	baseURL := os.Getenv("X402_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.example.com" // default for demonstration
	}

	// Create schemas using helper functions
	schemaProvider := schema.NewPathBased(map[string]*x402.EndpointSchema{
		"/twitter/user/followers": createTwitterFollowersSchema(),
		"/api/data":               createDataSubmitSchema(),
		"/api/weather":            createWeatherSchema(),
	}, nil)

	// Create resource provider with custom URLs and descriptions
	resourceProvider := resource.NewPathBased(map[string]*resource.Metadata{
		"/twitter/user/followers": {
			URL:         baseURL + "/twitter/user/followers",
			Description: "Get Twitter user followers with pagination support",
		},
		"/api/data": {
			URL:         baseURL + "/api/data",
			Description: "Submit and process data with optional filters",
		},
		"/api/weather": {
			URL:         baseURL + "/api/weather",
			Description: "Get current weather information for a city",
		},
	}, nil, baseURL)

	return &localx402.Config{
		RecipientAddress: recipientAddress,
		Network:          network,
		FacilitatorURL:   facilitatorURL,
		PricingStrategy:  pricing.NewFixed(decimal.RequireFromString("0.3")),
		SchemaProvider:   schemaProvider,   // Add schema provider
		ResourceProvider: resourceProvider, // Add resource provider
	}, nil
}

func setupRoutes(r *gin.Engine, config *localx402.Config) {
	// Free endpoint - no payment required
	r.GET("/health", func(c *gin.Context) {
		c.JSON(httpStatusOK, gin.H{
			"status": "ok",
		})
	})

	// Protected endpoints with schema
	protected := r.Group("/")
	protected.Use(ginx402.NewMiddleware(config))
	{
		// Twitter followers endpoint (GET with query params)
		protected.GET("/twitter/user/followers", func(c *gin.Context) {
			userName := c.Query("userName")
			cursor := c.DefaultQuery("cursor", "")
			pageSize := c.DefaultQuery("pageSize", "200")

			c.JSON(httpStatusOK, gin.H{
				"message": "Twitter followers data",
				"data": gin.H{
					"userName": userName,
					"cursor":   cursor,
					"pageSize": pageSize,
					"followers": []gin.H{
						{"id": "123", "username": "user1", "followedAt": "2024-01-01T00:00:00Z"},
						{"id": "456", "username": "user2", "followedAt": "2024-01-02T00:00:00Z"},
					},
				},
			})
		})

		// Data submission endpoint (POST with body)
		protected.POST("/api/data", func(c *gin.Context) {
			var body map[string]any
			if err := c.BindJSON(&body); err != nil {
				c.JSON(httpStatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			c.JSON(httpStatusOK, gin.H{
				"message": "Data processed successfully",
				"results": []string{"result1", "result2"},
				"total":   defaultResultCount,
			})
		})

		// Weather endpoint (GET with enum query param)
		protected.GET("/api/weather", func(c *gin.Context) {
			city := c.Query("city")
			units := c.DefaultQuery("units", "metric")

			c.JSON(httpStatusOK, gin.H{
				"city":        city,
				"temperature": sampleTemperature,
				"units":       units,
				"description": "Sunny",
			})
		})
	}
}

func main() {
	r := gin.Default()

	// Get configuration with schema support
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	setupRoutes(r, config)

	log.Print("Starting server on :8080")
	log.Printf("Network: %s", config.Network)
	log.Printf("Recipient: %s", config.RecipientAddress)
	log.Printf("Facilitator: %s", config.FacilitatorURL)
	log.Print("Schema support enabled for:")
	log.Print("  - /twitter/user/followers (GET)")
	log.Print("  - /api/data (POST)")
	log.Print("  - /api/weather (GET)")

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
