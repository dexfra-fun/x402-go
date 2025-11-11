package pricing

import (
	"context"

	"github.com/dexfra-fun/x402-go/pkg/x402"
)

// Fixed implements a fixed price strategy
type Fixed struct {
	amount float64
}

// NewFixed creates a new fixed pricing strategy
func NewFixed(amount float64) *Fixed {
	return &Fixed{amount: amount}
}

// GetPrice returns the fixed price for any resource
func (p *Fixed) GetPrice(ctx context.Context, resource x402.Resource) (float64, error) {
	return p.amount, nil
}

// Dynamic implements a dynamic pricing strategy using a custom fetcher
type Dynamic struct {
	fetcher PriceFetcher
}

// PriceFetcher defines the interface for fetching prices dynamically
type PriceFetcher interface {
	FetchPrice(ctx context.Context, resource x402.Resource) (float64, error)
}

// NewDynamic creates a new dynamic pricing strategy
func NewDynamic(fetcher PriceFetcher) *Dynamic {
	return &Dynamic{fetcher: fetcher}
}

// GetPrice fetches the price dynamically using the provided fetcher
func (p *Dynamic) GetPrice(ctx context.Context, resource x402.Resource) (float64, error) {
	return p.fetcher.FetchPrice(ctx, resource)
}

// PathBased implements pricing based on path patterns
type PathBased struct {
	prices       map[string]float64 // path -> price mapping
	defaultPrice float64
}

// NewPathBased creates a new path-based pricing strategy
func NewPathBased(prices map[string]float64, defaultPrice float64) *PathBased {
	return &PathBased{
		prices:       prices,
		defaultPrice: defaultPrice,
	}
}

// GetPrice returns price based on the resource path
func (p *PathBased) GetPrice(ctx context.Context, resource x402.Resource) (float64, error) {
	if price, ok := p.prices[resource.Path]; ok {
		return price, nil
	}
	return p.defaultPrice, nil
}

// MethodBased implements pricing based on HTTP methods
type MethodBased struct {
	prices       map[string]float64 // method -> price mapping
	defaultPrice float64
}

// NewMethodBased creates a new method-based pricing strategy
func NewMethodBased(prices map[string]float64, defaultPrice float64) *MethodBased {
	return &MethodBased{
		prices:       prices,
		defaultPrice: defaultPrice,
	}
}

// GetPrice returns price based on the HTTP method
func (p *MethodBased) GetPrice(ctx context.Context, resource x402.Resource) (float64, error) {
	if price, ok := p.prices[resource.Method]; ok {
		return price, nil
	}
	return p.defaultPrice, nil
}
