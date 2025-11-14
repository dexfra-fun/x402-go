# x402-go

[![Go Reference](https://pkg.go.dev/badge/github.com/dexfra-fun/x402-go.svg)](https://pkg.go.dev/github.com/dexfra-fun/x402-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/dexfra-fun/x402-go)](https://goreportcard.com/report/github.com/dexfra-fun/x402-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go middleware library for implementing the [x402 payment protocol](https://github.com/coinbase/x402) in your APIs. Built for the Dexfra ecosystem but usable by any Go HTTP server.

## Features

- 🔌 **Framework Support**: Gin, Chi, Fiber, and standard `net/http`
- 💰 **Flexible Pricing**: Fixed or dynamic pricing strategies
- 📋 **Schema Support**: Define input/output schemas for API endpoints
- ⚡ **Facilitator Integration**: Built-in support for x402 facilitators
- 🚀 **Performance**: Fee payer caching to reduce facilitator calls
- 🌐 **Multi-Chain**: Support for Solana, EVM, and other networks
- 🧪 **Well Tested**: Comprehensive unit and integration tests

## Quick Start

### Installation

```bash
go get github.com/dexfra-fun/x402-go
```

### Basic Usage with Gin

```go
package main

import (
    "github.com/dexfra-fun/x402-go/pkg/x402"
    ginx402 "github.com/dexfra-fun/x402-go/pkg/adapters/gin"
    "github.com/dexfra-fun/x402-go/pkg/pricing"
    "github.com/gin-gonic/gin"
    "github.com/shopspring/decimal"
)

func main() {
    r := gin.Default()
    
    // Configure x402 middleware
    config := &x402.Config{
        RecipientAddress: "YOUR_WALLET_ADDRESS",
        Network:          "solana-devnet",
        FacilitatorURL:   "https://facilitator.payai.network",
        PricingStrategy:  pricing.NewFixed(decimal.RequireFromString("0.001")), // 0.001 USDC
    }
    
    // Apply middleware to protected routes
    r.GET("/api/data", ginx402.NewMiddleware(config), func(c *gin.Context) {
        c.JSON(200, gin.H{"data": "protected content"})
    })
    
    r.Run(":8080")
}
```

### Usage with Chi

```go
package main

import (
    "net/http"
    
    "github.com/dexfra-fun/x402-go/pkg/x402"
    chix402 "github.com/dexfra-fun/x402-go/pkg/adapters/chi"
    "github.com/dexfra-fun/x402-go/pkg/pricing"
    "github.com/go-chi/chi/v5"
)

func main() {
    r := chi.NewRouter()
    
    config := &x402.Config{
        RecipientAddress: "YOUR_WALLET_ADDRESS",
        Network:          "solana-devnet",
        FacilitatorURL:   "https://facilitator.payai.network",
        PricingStrategy:  pricing.NewFixed(0.001),
    }
    
    // Apply middleware
    r.Use(chix402.NewMiddleware(config))
    
    r.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`{"data": "protected content"}`))
    })
    
    http.ListenAndServe(":8080", r)
}
```

### Usage with Standard HTTP

```go
package main

import (
    "net/http"
    
    "github.com/dexfra-fun/x402-go/pkg/x402"
    httpx402 "github.com/dexfra-fun/x402-go/pkg/adapters/http"
    "github.com/dexfra-fun/x402-go/pkg/pricing"
)

func main() {
    config := &x402.Config{
        RecipientAddress: "YOUR_WALLET_ADDRESS",
        Network:          "solana-devnet",
        FacilitatorURL:   "https://facilitator.payai.network",
        PricingStrategy:  pricing.NewFixed(0.001),
    }
    
    mux := http.NewServeMux()
    mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`{"data": "protected content"}`))
    })
    
    // Wrap with x402 middleware
    handler := httpx402.NewMiddleware(config)(mux)
    
    http.ListenAndServe(":8080", handler)
}
```

## Dynamic Pricing

Implement custom pricing logic:

```go
import "github.com/shopspring/decimal"

type DatabasePricer struct {
    db *sql.DB
}

func (p *DatabasePricer) GetPrice(ctx context.Context, resource x402.Resource) (decimal.Decimal, error) {
    var priceStr string
    err := p.db.QueryRowContext(ctx, 
        "SELECT price FROM api_pricing WHERE path = ? AND method = ?",
        resource.Path, resource.Method,
    ).Scan(&priceStr)
    if err != nil {
        return decimal.Zero, err
    }
    return decimal.NewFromString(priceStr)
}

// Use it
config := &x402.Config{
    RecipientAddress: "YOUR_WALLET_ADDRESS",
    Network:          "solana-devnet",
    FacilitatorURL:   "https://facilitator.payai.network",
    PricingStrategy:  &DatabasePricer{db: db},
}
```

## Schema Support

Define input/output schemas for your API endpoints according to the [x402 specification](https://github.com/coinbase/x402). Schemas are automatically included in 402 responses to help clients understand your API structure.

### Basic Schema Usage

```go
import (
    "github.com/dexfra-fun/x402-go"
    "github.com/dexfra-fun/x402-go/pkg/schema"
)

// Define schema for an endpoint
followersSchema := &x402.EndpointSchema{
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
                Description: "pagination cursor",
            },
        },
        HeaderFields: map[string]*x402.FieldDef{
            "aisa-payment": {
                Type:        "string",
                Required:    false,
                Description: "Payment proof header",
            },
        },
    },
}

// Use path-based schema provider
schemaProvider := schema.NewPathBased(map[string]*x402.EndpointSchema{
    "/twitter/user/followers": followersSchema,
}, nil)

config := &x402.Config{
    RecipientAddress: "YOUR_WALLET_ADDRESS",
    Network:          "solana-devnet",
    FacilitatorURL:   "https://facilitator.payai.network",
    PricingStrategy:  pricing.NewFixed(decimal.NewFromFloat(0.001)),
    SchemaProvider:   schemaProvider, // Add schema provider
}
```

### Schema Builder API

Use the fluent builder API for easier schema construction:

```go
schema := schema.NewEndpointSchema().
    WithInput(
        schema.NewInputSchema("POST").
            WithBodyType("json").
            WithBodyField("query", schema.NewFieldDef("string", true, "Search query")).
            WithBodyField("limit", schema.NewFieldDef("integer", false, "Result limit")).
            WithHeaderField("X-API-Key", schema.NewFieldDef("string", false, "API key")).
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
```

### Schema Provider Strategies

**Static Schema** - Same schema for all endpoints:
```go
schema := schema.NewStatic(mySchema)
```

**Path-Based Schema** - Different schemas per path:
```go
schema := schema.NewPathBased(map[string]*x402.EndpointSchema{
    "/api/users":    usersSchema,
    "/api/products": productsSchema,
}, defaultSchema)
```

**Dynamic Schema** - Custom logic:
```go
type MySchemaFetcher struct {
    db *sql.DB
}

func (f *MySchemaFetcher) FetchSchema(ctx context.Context, resource x402.Resource) (*x402.EndpointSchema, error) {
    // Fetch from database, external service, etc.
    return schemaFromDB, nil
}

schema := schema.NewDynamic(&MySchemaFetcher{db: db})
```

### Advanced Field Definitions

**Nested Objects:**
```go
schema.NewObjectField(map[string]*x402.FieldDef{
    "city":    schema.NewFieldDef("string", true, "City name"),
    "country": schema.NewFieldDef("string", true, "Country code"),
}, true, "Address object")
```

**Enum Fields:**
```go
schema.NewEnumField([]string{"active", "inactive", "pending"}, false, "Status value")
```

**Conditional Requirements:**
```go
schema.NewConditionalField("string", []string{"otherField"}, "Required when otherField is present")
```

## Configuration Options

```go
type Config struct {
    // Required
    RecipientAddress string          // Your wallet address for receiving payments
    Network          string          // e.g., "solana-devnet", "solana-mainnet"
    FacilitatorURL   string          // Facilitator service URL
    PricingStrategy  PricingStrategy // How to price API calls
    
    // Optional
    CacheTTL         time.Duration   // Fee payer cache duration (default: 5 minutes)
    Networks         map[string]NetworkConfig // Custom network configurations
    Logger           Logger          // Custom logger
}
```

## Supported Networks

- Solana (devnet, mainnet)
- EVM chains (coming soon)
- Custom network configurations

## Examples

See the [examples](./examples) directory for complete working examples:

- [Gin example](./examples/gin/main.go)
- [Chi example](./examples/chi/main.go)
- [Standard HTTP example](./examples/http/main.go)
- [Schema support example](./examples/schema/main.go) - Demonstrates schema definition and usage

## Documentation

- [API Reference](https://pkg.go.dev/github.com/dexfra-fun/x402-go)
- [x402 Protocol Specification](https://github.com/coinbase/x402)
- [Dexfra Documentation](https://docs.dexfra.fun)

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built on top of [mark3labs/x402-go](https://github.com/mark3labs/x402-go)
- Part of the [Dexfra](https://dexfra.fun) ecosystem
- Implements the [Coinbase x402 specification](https://github.com/coinbase/x402)
