package schema

import (
	"context"

	x402 "github.com/dexfra-fun/x402-go"
	localx402 "github.com/dexfra-fun/x402-go/pkg/x402"
)

// Static provides a static schema for all endpoints.
// This is useful when all endpoints share the same schema.
type Static struct {
	schema *x402.EndpointSchema
}

// NewStatic creates a new static schema provider.
func NewStatic(schema *x402.EndpointSchema) *Static {
	return &Static{
		schema: schema,
	}
}

// GetSchema returns the static schema for any resource.
func (s *Static) GetSchema(_ context.Context, _ localx402.Resource) (*x402.EndpointSchema, error) {
	return s.schema, nil
}