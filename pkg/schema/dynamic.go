package schema

import (
	"context"

	x402 "github.com/dexfra-fun/x402-go"
	localx402 "github.com/dexfra-fun/x402-go/pkg/x402"
)

// SchemaFetcher defines the interface for fetching schemas dynamically.
// Users can implement this interface to fetch schemas from databases,
// external services, or any custom logic.
type SchemaFetcher interface {
	FetchSchema(ctx context.Context, resource localx402.Resource) (*x402.EndpointSchema, error)
}

// Dynamic provides schemas using a custom fetcher.
// This allows for complex schema resolution logic like database lookups,
// API calls, or other dynamic sources.
type Dynamic struct {
	fetcher SchemaFetcher
}

// NewDynamic creates a new dynamic schema provider.
func NewDynamic(fetcher SchemaFetcher) *Dynamic {
	return &Dynamic{
		fetcher: fetcher,
	}
}

// GetSchema fetches the schema dynamically using the provided fetcher.
func (d *Dynamic) GetSchema(ctx context.Context, resource localx402.Resource) (*x402.EndpointSchema, error) {
	return d.fetcher.FetchSchema(ctx, resource)
}