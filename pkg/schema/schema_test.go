package schema

import (
	"testing"

	x402 "github.com/dexfra-fun/x402-go"
)

func TestNewFieldDef(t *testing.T) {
	field := NewFieldDef("string", true, "Test field")

	if field.Type != "string" {
		t.Errorf("Expected type 'string', got '%s'", field.Type)
	}
	if field.Required != true {
		t.Errorf("Expected required to be true")
	}
	if field.Description != "Test field" {
		t.Errorf("Expected description 'Test field', got '%s'", field.Description)
	}
}

func TestNewEnumField(t *testing.T) {
	enum := []string{"active", "inactive", "pending"}
	field := NewEnumField(enum, false, "Status field")

	if field.Type != "string" {
		t.Errorf("Expected type 'string', got '%s'", field.Type)
	}
	if len(field.Enum) != 3 {
		t.Errorf("Expected 3 enum values, got %d", len(field.Enum))
	}
	if field.Enum[0] != "active" {
		t.Errorf("Expected first enum value 'active', got '%s'", field.Enum[0])
	}
}

func TestNewObjectField(t *testing.T) {
	properties := map[string]*x402.FieldDef{
		"name": NewFieldDef("string", true, "Name"),
		"age":  NewFieldDef("integer", false, "Age"),
	}

	field := NewObjectField(properties, true, "User object")

	if field.Type != "object" {
		t.Errorf("Expected type 'object', got '%s'", field.Type)
	}
	if len(field.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(field.Properties))
	}
	if field.Properties["name"].Type != "string" {
		t.Errorf("Expected name property type 'string', got '%s'", field.Properties["name"].Type)
	}
}

func TestNewConditionalField(t *testing.T) {
	requiredWhen := []string{"otherField", "anotherField"}
	field := NewConditionalField("string", requiredWhen, "Conditional field")

	if field.Type != "string" {
		t.Errorf("Expected type 'string', got '%s'", field.Type)
	}

	// Required should be a slice
	if required, ok := field.Required.([]string); ok {
		if len(required) != 2 {
			t.Errorf("Expected 2 conditional requirements, got %d", len(required))
		}
	} else {
		t.Errorf("Expected Required to be []string, got %T", field.Required)
	}
}

func TestInputSchemaBuilder(t *testing.T) {
	schema := NewInputSchema("POST").
		WithBodyType("json").
		WithQueryParam("page", NewFieldDef("integer", false, "Page number")).
		WithBodyField("name", NewFieldDef("string", true, "User name")).
		WithHeaderField("X-API-Key", NewFieldDef("string", false, "API key")).
		Build()

	if schema.Type != "http" {
		t.Errorf("Expected type 'http', got '%s'", schema.Type)
	}
	if schema.Method != "POST" {
		t.Errorf("Expected method 'POST', got '%s'", schema.Method)
	}
	if schema.BodyType != "json" {
		t.Errorf("Expected bodyType 'json', got '%s'", schema.BodyType)
	}
	if len(schema.QueryParams) != 1 {
		t.Errorf("Expected 1 query param, got %d", len(schema.QueryParams))
	}
	if len(schema.BodyFields) != 1 {
		t.Errorf("Expected 1 body field, got %d", len(schema.BodyFields))
	}
	if len(schema.HeaderFields) != 1 {
		t.Errorf("Expected 1 header field, got %d", len(schema.HeaderFields))
	}
}

func TestEndpointSchemaBuilder(t *testing.T) {
	input := NewInputSchema("GET").Build()
	output := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"result": map[string]any{"type": "string"},
		},
	}

	schema := NewEndpointSchema().
		WithInput(input).
		WithOutput(output).
		Build()

	if schema.Input == nil {
		t.Error("Expected input schema to be set")
	}
	if schema.Output == nil {
		t.Error("Expected output schema to be set")
	}
	if schema.Output["type"] != "object" {
		t.Errorf("Expected output type 'object', got '%v'", schema.Output["type"])
	}
}