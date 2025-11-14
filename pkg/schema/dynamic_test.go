package schema

import (
	"context"
	"errors"
	"testing"

	x402 "github.com/dexfra-fun/x402-go"
	localx402 "github.com/dexfra-fun/x402-go/pkg/x402"
)

// MockSchemaFetcher is a mock implementation of SchemaFetcher for testing
type MockSchemaFetcher struct {
	schemas map[string]*x402.EndpointSchema
	err     error
}

func (m *MockSchemaFetcher) FetchSchema(_ context.Context, resource localx402.Resource) (*x402.EndpointSchema, error) {
	if m.err != nil {
		return nil, m.err
	}
	if schema, ok := m.schemas[resource.Path]; ok {
		return schema, nil
	}
	return nil, nil
}

func TestDynamic(t *testing.T) {
	// Create test schemas
	usersSchema := &x402.EndpointSchema{
		Input: &x402.InputSchema{
			Type:   "http",
			Method: "GET",
		},
	}

	productsSchema := &x402.EndpointSchema{
		Input: &x402.InputSchema{
			Type:   "http",
			Method: "POST",
		},
	}

	// Create mock fetcher
	fetcher := &MockSchemaFetcher{
		schemas: map[string]*x402.EndpointSchema{
			"/api/users":    usersSchema,
			"/api/products": productsSchema,
		},
	}

	provider := NewDynamic(fetcher)

	tests := []struct {
		name           string
		path           string
		expectedSchema *x402.EndpointSchema
	}{
		{
			name:           "fetch users schema",
			path:           "/api/users",
			expectedSchema: usersSchema,
		},
		{
			name:           "fetch products schema",
			path:           "/api/products",
			expectedSchema: productsSchema,
		},
		{
			name:           "no schema for unknown path",
			path:           "/api/unknown",
			expectedSchema: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := localx402.Resource{Path: tt.path, Method: "GET"}
			schema, err := provider.GetSchema(context.Background(), resource)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if schema != tt.expectedSchema {
				t.Error("Expected specific schema instance")
			}
		})
	}
}

func TestDynamicError(t *testing.T) {
	expectedErr := errors.New("database connection failed")
	fetcher := &MockSchemaFetcher{
		err: expectedErr,
	}

	provider := NewDynamic(fetcher)

	resource := localx402.Resource{Path: "/api/test", Method: "GET"}
	schema, err := provider.GetSchema(context.Background(), resource)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err != expectedErr {
		t.Errorf("Expected error '%v', got '%v'", expectedErr, err)
	}
	if schema != nil {
		t.Error("Expected nil schema on error")
	}
}

func TestDynamicWithContext(t *testing.T) {
	// Test that context is properly passed to fetcher
	fetcher := &MockSchemaFetcher{
		schemas: map[string]*x402.EndpointSchema{
			"/api/test": {
				Input: &x402.InputSchema{
					Type:   "http",
					Method: "GET",
				},
			},
		},
	}

	provider := NewDynamic(fetcher)

	// Use a canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	resource := localx402.Resource{Path: "/api/test", Method: "GET"}
	// Note: Our mock doesn't check context cancellation, but in real implementation it would
	schema, err := provider.GetSchema(ctx, resource)

	// With our simple mock, it should still work
	if err != nil {
		t.Errorf("Mock doesn't respect context cancellation, got error: %v", err)
	}
	if schema == nil {
		t.Error("Expected schema even with canceled context in mock")
	}
}