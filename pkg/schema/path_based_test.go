package schema

import (
	"context"
	"net/http"
	"testing"

	x402 "github.com/dexfra-fun/x402-go"
	localx402 "github.com/dexfra-fun/x402-go/pkg/x402"
)

func testPathBasedExactMatch(t *testing.T, provider *PathBased, usersSchema, productsSchema *x402.EndpointSchema) {
	t.Helper()
	tests := []struct {
		name           string
		path           string
		expectedMethod string
		expectedSchema *x402.EndpointSchema
	}{
		{
			name:           "exact match users",
			path:           "/api/users",
			expectedMethod: http.MethodGet,
			expectedSchema: usersSchema,
		},
		{
			name:           "exact match products",
			path:           "/api/products",
			expectedMethod: http.MethodPost,
			expectedSchema: productsSchema,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := localx402.Resource{Path: tt.path, Method: http.MethodGet}
			schema, err := provider.GetSchema(context.Background(), resource)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if schema != tt.expectedSchema {
				t.Error("Expected specific schema instance, got different one")
			}
			if schema.Input.Method != tt.expectedMethod {
				t.Errorf("Expected method '%s', got '%s'", tt.expectedMethod, schema.Input.Method)
			}
		})
	}
}

func TestPathBased(t *testing.T) {
	// Create test schemas
	usersSchema := &x402.EndpointSchema{
		Input: &x402.InputSchema{
			Type:   "http",
			Method: http.MethodGet,
			QueryParams: map[string]*x402.FieldDef{
				"page": NewFieldDef("integer", false, "Page number"),
			},
		},
	}

	productsSchema := &x402.EndpointSchema{
		Input: &x402.InputSchema{
			Type:   "http",
			Method: http.MethodPost,
			BodyFields: map[string]*x402.FieldDef{
				"name": NewFieldDef("string", true, "Product name"),
			},
		},
	}

	defaultSchema := &x402.EndpointSchema{
		Input: &x402.InputSchema{
			Type:   "http",
			Method: http.MethodGet,
		},
	}

	// Create path-based provider
	provider := NewPathBased(map[string]*x402.EndpointSchema{
		"/api/users":    usersSchema,
		"/api/products": productsSchema,
	}, defaultSchema)

	// Test exact matches
	testPathBasedExactMatch(t, provider, usersSchema, productsSchema)

	// Test default fallback
	t.Run("no match returns default", func(t *testing.T) {
		resource := localx402.Resource{Path: "/api/other", Method: http.MethodGet}
		schema, err := provider.GetSchema(context.Background(), resource)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if schema != defaultSchema {
			t.Error("Expected default schema")
		}
		if schema.Input.Method != http.MethodGet {
			t.Errorf("Expected method 'GET', got '%s'", schema.Input.Method)
		}
	})
}

func TestPathBasedNoDefault(t *testing.T) {
	schema := &x402.EndpointSchema{
		Input: &x402.InputSchema{
			Type:   "http",
			Method: http.MethodGet,
		},
	}

	provider := NewPathBased(map[string]*x402.EndpointSchema{
		"/api/users": schema,
	}, nil)

	// Test with non-matching path
	resource := localx402.Resource{Path: "/api/other", Method: http.MethodGet}
	result, err := provider.GetSchema(context.Background(), resource)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != nil {
		t.Error("Expected nil schema when no match and no default")
	}
}

func TestPathBasedAddSchema(t *testing.T) {
	provider := NewPathBased(nil, nil)

	// Initially, should return nil for any path
	resource := localx402.Resource{Path: "/api/test", Method: http.MethodGet}
	schema, err := provider.GetSchema(context.Background(), resource)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if schema != nil {
		t.Error("Expected nil schema before adding")
	}

	// Add a schema
	testSchema := &x402.EndpointSchema{
		Input: &x402.InputSchema{
			Type:   "http",
			Method: http.MethodPost,
		},
	}
	provider.AddSchema("/api/test", testSchema)

	// Now should return the added schema
	schema, err = provider.GetSchema(context.Background(), resource)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if schema != testSchema {
		t.Error("Expected added schema to be returned")
	}
}

func TestPathBasedSetDefaultSchema(t *testing.T) {
	provider := NewPathBased(nil, nil)

	defaultSchema := &x402.EndpointSchema{
		Input: &x402.InputSchema{
			Type:   "http",
			Method: http.MethodGet,
		},
	}

	provider.SetDefaultSchema(defaultSchema)

	// Any path should now return the default schema
	resource := localx402.Resource{Path: "/any/path", Method: http.MethodGet}
	schema, err := provider.GetSchema(context.Background(), resource)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if schema != defaultSchema {
		t.Error("Expected default schema to be returned")
	}
}