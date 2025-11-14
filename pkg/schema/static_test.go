package schema

import (
	"context"
	"testing"

	x402 "github.com/dexfra-fun/x402-go"
	localx402 "github.com/dexfra-fun/x402-go/pkg/x402"
)

func TestStatic(t *testing.T) {
	// Create a test schema
	testSchema := &x402.EndpointSchema{
		Input: &x402.InputSchema{
			Type:   "http",
			Method: "GET",
			QueryParams: map[string]*x402.FieldDef{
				"id": {
					Type:        "string",
					Required:    true,
					Description: "Item ID",
				},
			},
		},
	}

	// Create static provider
	provider := NewStatic(testSchema)

	// Test with different resources
	resources := []localx402.Resource{
		{Path: "/api/users", Method: "GET"},
		{Path: "/api/products", Method: "POST"},
		{Path: "/different/path", Method: "DELETE"},
	}

	for _, resource := range resources {
		t.Run(resource.Path+"_"+resource.Method, func(t *testing.T) {
			schema, err := provider.GetSchema(context.Background(), resource)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if schema == nil {
				t.Error("Expected schema, got nil")
			}
			if schema.Input.Method != "GET" {
				t.Errorf("Expected method 'GET', got '%s'", schema.Input.Method)
			}
			// Verify it's the same schema instance
			if schema != testSchema {
				t.Error("Expected same schema instance for all resources")
			}
		})
	}
}

func TestStaticNilSchema(t *testing.T) {
	provider := NewStatic(nil)

	resource := localx402.Resource{Path: "/api/test", Method: "GET"}
	schema, err := provider.GetSchema(context.Background(), resource)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if schema != nil {
		t.Error("Expected nil schema")
	}
}