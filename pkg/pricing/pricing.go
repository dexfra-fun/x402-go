// Package pricing provides various pricing strategies for x402 payment middleware.
package pricing

import (
	"context"

	"github.com/dexfra-fun/x402-go/pkg/x402"
	"github.com/shopspring/decimal"
)

// Fixed implements a fixed price strategy.
type Fixed struct {
	amount decimal.Decimal
}

// NewFixed creates a new fixed pricing strategy.
func NewFixed(amount decimal.Decimal) *Fixed {
	return &Fixed{amount: amount}
}

// NewFixedFromFloat creates a new fixed pricing strategy from float64.
func NewFixedFromFloat(amount float64) *Fixed {
	return &Fixed{amount: decimal.NewFromFloat(amount)}
}

// NewFixedFromString creates a new fixed pricing strategy from string.
func NewFixedFromString(amount string) (*Fixed, error) {
	amt, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, err
	}
	return &Fixed{amount: amt}, nil
}

// GetPrice returns the fixed price for any resource.
func (p *Fixed) GetPrice(_ context.Context, _ x402.Resource) (decimal.Decimal, error) {
	return p.amount, nil
}

// Dynamic implements a dynamic pricing strategy using a custom fetcher.
type Dynamic struct {
	fetcher PriceFetcher
}

// PriceFetcher defines the interface for fetching prices dynamically.
type PriceFetcher interface {
	FetchPrice(ctx context.Context, resource x402.Resource) (decimal.Decimal, error)
}

// NewDynamic creates a new dynamic pricing strategy.
func NewDynamic(fetcher PriceFetcher) *Dynamic {
	return &Dynamic{fetcher: fetcher}
}

// GetPrice fetches the price dynamically using the provided fetcher.
func (p *Dynamic) GetPrice(ctx context.Context, resource x402.Resource) (decimal.Decimal, error) {
	return p.fetcher.FetchPrice(ctx, resource)
}

// PathBased implements pricing based on path patterns.
type PathBased struct {
	prices       map[string]decimal.Decimal // path -> price mapping
	defaultPrice decimal.Decimal
}

// NewPathBased creates a new path-based pricing strategy.
func NewPathBased(prices map[string]decimal.Decimal, defaultPrice decimal.Decimal) *PathBased {
	return &PathBased{
		prices:       prices,
		defaultPrice: defaultPrice,
	}
}

// NewPathBasedFromFloat creates path-based pricing from float64 map.
func NewPathBasedFromFloat(prices map[string]float64, defaultPrice float64) *PathBased {
	decPrices := make(map[string]decimal.Decimal, len(prices))
	for path, price := range prices {
		decPrices[path] = decimal.NewFromFloat(price)
	}
	return &PathBased{
		prices:       decPrices,
		defaultPrice: decimal.NewFromFloat(defaultPrice),
	}
}

// GetPrice returns price based on the resource path.
func (p *PathBased) GetPrice(_ context.Context, resource x402.Resource) (decimal.Decimal, error) {
	if price, ok := p.prices[resource.Path]; ok {
		return price, nil
	}
	return p.defaultPrice, nil
}

// MethodBased implements pricing based on HTTP methods.
type MethodBased struct {
	prices       map[string]decimal.Decimal // method -> price mapping
	defaultPrice decimal.Decimal
}

// NewMethodBased creates a new method-based pricing strategy.
func NewMethodBased(prices map[string]decimal.Decimal, defaultPrice decimal.Decimal) *MethodBased {
	return &MethodBased{
		prices:       prices,
		defaultPrice: defaultPrice,
	}
}

// NewMethodBasedFromFloat creates method-based pricing from float64 map.
func NewMethodBasedFromFloat(prices map[string]float64, defaultPrice float64) *MethodBased {
	decPrices := make(map[string]decimal.Decimal, len(prices))
	for method, price := range prices {
		decPrices[method] = decimal.NewFromFloat(price)
	}
	return &MethodBased{
		prices:       decPrices,
		defaultPrice: decimal.NewFromFloat(defaultPrice),
	}
}

// GetPrice returns price based on the HTTP method.
func (p *MethodBased) GetPrice(
	_ context.Context,
	resource x402.Resource,
) (decimal.Decimal, error) {
	if price, ok := p.prices[resource.Method]; ok {
		return price, nil
	}
	return p.defaultPrice, nil
}
