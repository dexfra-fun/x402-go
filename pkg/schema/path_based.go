package schema

import (
	"context"

	x402 "github.com/dexfra-fun/x402-go"
	localx402 "github.com/dexfra-fun/x402-go/pkg/x402"
)

// PathBased provides different schemas based on the request path.
// This allows different endpoints to have different schema definitions.
type PathBased struct {
	schemas       map[string]*x402.EndpointSchema
	defaultSchema *x402.EndpointSchema
}

// NewPathBased creates a new path-based schema provider.
// The schemas map keys should be the path patterns to match.
// If defaultSchema is provided, it will be used when no path matches.
func NewPathBased(schemas map[string]*x402.EndpointSchema, defaultSchema *x402.EndpointSchema) *PathBased {
	return &PathBased{
		schemas:       schemas,
		defaultSchema: defaultSchema,
	}
}

// GetSchema returns the schema for the given resource path.
// If no matching schema is found, returns the default schema (if configured).
func (p *PathBased) GetSchema(_ context.Context, resource localx402.Resource) (*x402.EndpointSchema, error) {
	// Try exact path match first
	if schema, ok := p.schemas[resource.Path]; ok {
		return schema, nil
	}

	// Return default schema if configured
	return p.defaultSchema, nil
}

// AddSchema adds or updates a schema for a specific path.
func (p *PathBased) AddSchema(path string, schema *x402.EndpointSchema) {
	if p.schemas == nil {
		p.schemas = make(map[string]*x402.EndpointSchema)
	}
	p.schemas[path] = schema
}

// SetDefaultSchema sets the default schema to use when no path matches.
func (p *PathBased) SetDefaultSchema(schema *x402.EndpointSchema) {
	p.defaultSchema = schema
}