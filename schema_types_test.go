package x402

import (
	"encoding/json"
	"testing"
)

func TestInputSchema_Discoverable_DefaultTrue(t *testing.T) {
	tests := []struct {
		name     string
		schema   *InputSchema
		expected bool
	}{
		{
			name: "not set - defaults to true",
			schema: &InputSchema{
				Type:   "http",
				Method: "GET",
			},
			expected: true,
		},
		{
			name: "explicitly set to true",
			schema: &InputSchema{
				Type:         "http",
				Method:       "POST",
				Discoverable: boolPtr(true),
			},
			expected: true,
		},
		{
			name: "explicitly set to false",
			schema: &InputSchema{
				Type:         "http",
				Method:       "DELETE",
				Discoverable: boolPtr(false),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test IsDiscoverable method
			if got := tt.schema.IsDiscoverable(); got != tt.expected {
				t.Errorf("IsDiscoverable() = %v, want %v", got, tt.expected)
			}

			// Test JSON marshaling
			data, err := json.Marshal(tt.schema)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			var result map[string]any
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			discoverable, ok := result["discoverable"].(bool)
			if !ok {
				t.Fatal("discoverable field not found or not a bool in JSON output")
			}

			if discoverable != tt.expected {
				t.Errorf("JSON discoverable = %v, want %v", discoverable, tt.expected)
			}
		})
	}
}

// boolPtr returns a pointer to a bool value.
func boolPtr(b bool) *bool {
	return &b
}