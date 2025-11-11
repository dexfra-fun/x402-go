# x402-go

[![Go Reference](https://pkg.go.dev/badge/github.com/dexfra-fun/x402-go.svg)](https://pkg.go.dev/github.com/dexfra-fun/x402-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/dexfra-fun/x402-go)](https://goreportcard.com/report/github.com/dexfra-fun/x402-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go middleware library for implementing the [x402 payment protocol](https://github.com/coinbase/x402) in your APIs. Built for the Dexfra ecosystem but usable by any Go HTTP server.

## Features

- 🔌 **Framework Support**: Gin, Chi, and standard `net/http`
- 💰 **Flexible Pricing**: Fixed or dynamic pricing strategies
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
)

func main() {
    r := gin.Default()
    
    // Configure x402 middleware
    config := &x402.Config{
        RecipientAddress: "YOUR_WALLET_ADDRESS",
        Network:          "solana-devnet",
        FacilitatorURL:   "https://facilitator.payai.network",
        PricingStrategy:  pricing.NewFixed(0.001), // 0.001 USDC per call
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
type DatabasePricer struct {
    db *sql.DB
}

func (p *DatabasePricer) GetPrice(ctx context.Context, resource x402.Resource) (float64, error) {
    var price float64
    err := p.db.QueryRowContext(ctx, 
        "SELECT price FROM api_pricing WHERE path = ? AND method = ?",
        resource.Path, resource.Method,
    ).Scan(&price)
    return price, err
}

// Use it
config := &x402.Config{
    RecipientAddress: "YOUR_WALLET_ADDRESS",
    Network:          "solana-devnet",
    FacilitatorURL:   "https://facilitator.payai.network",
    PricingStrategy:  &DatabasePricer{db: db},
}
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
