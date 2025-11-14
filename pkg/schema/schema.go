// Package schema provides implementations of SchemaProvider for x402 middleware.
// It offers various strategies for providing API endpoint schemas including static,
// path-based, and dynamic schema providers.
package schema

import (
	x402 "github.com/dexfra-fun/x402-go"
)

// NewFieldDef creates a basic field definition.
func NewFieldDef(fieldType string, required bool, description string) *x402.FieldDef {
	return &x402.FieldDef{
		Type:        fieldType,
		Required:    required,
		Description: description,
	}
}

// NewEnumField creates a field definition with enumerated values.
func NewEnumField(enum []string, required bool, description string) *x402.FieldDef {
	return &x402.FieldDef{
		Type:        "string",
		Enum:        enum,
		Required:    required,
		Description: description,
	}
}

// NewObjectField creates a field definition for nested objects.
func NewObjectField(properties map[string]*x402.FieldDef, required bool, description string) *x402.FieldDef {
	return &x402.FieldDef{
		Type:        "object",
		Properties:  properties,
		Required:    required,
		Description: description,
	}
}

// NewConditionalField creates a field with conditional requirements.
// The required parameter should be a slice of field names that make this field required.
func NewConditionalField(fieldType string, requiredWhen []string, description string) *x402.FieldDef {
	return &x402.FieldDef{
		Type:        fieldType,
		Required:    requiredWhen,
		Description: description,
	}
}

// InputSchemaBuilder provides a fluent API for building InputSchema.
type InputSchemaBuilder struct {
	schema *x402.InputSchema
}

// NewInputSchema creates a new InputSchemaBuilder with HTTP type and method.
func NewInputSchema(method string) *InputSchemaBuilder {
	return &InputSchemaBuilder{
		schema: &x402.InputSchema{
			Type:   "http",
			Method: method,
		},
	}
}

// WithBodyType sets the body content type.
func (b *InputSchemaBuilder) WithBodyType(bodyType string) *InputSchemaBuilder {
	b.schema.BodyType = bodyType
	return b
}

// WithQueryParam adds a query parameter to the schema.
func (b *InputSchemaBuilder) WithQueryParam(name string, field *x402.FieldDef) *InputSchemaBuilder {
	if b.schema.QueryParams == nil {
		b.schema.QueryParams = make(map[string]*x402.FieldDef)
	}
	b.schema.QueryParams[name] = field
	return b
}

// WithBodyField adds a body field to the schema.
func (b *InputSchemaBuilder) WithBodyField(name string, field *x402.FieldDef) *InputSchemaBuilder {
	if b.schema.BodyFields == nil {
		b.schema.BodyFields = make(map[string]*x402.FieldDef)
	}
	b.schema.BodyFields[name] = field
	return b
}

// WithHeaderField adds a header field to the schema.
func (b *InputSchemaBuilder) WithHeaderField(name string, field *x402.FieldDef) *InputSchemaBuilder {
	if b.schema.HeaderFields == nil {
		b.schema.HeaderFields = make(map[string]*x402.FieldDef)
	}
	b.schema.HeaderFields[name] = field
	return b
}

// Build returns the constructed InputSchema.
func (b *InputSchemaBuilder) Build() *x402.InputSchema {
	return b.schema
}

// EndpointSchemaBuilder provides a fluent API for building EndpointSchema.
type EndpointSchemaBuilder struct {
	schema *x402.EndpointSchema
}

// NewEndpointSchema creates a new EndpointSchemaBuilder.
func NewEndpointSchema() *EndpointSchemaBuilder {
	return &EndpointSchemaBuilder{
		schema: &x402.EndpointSchema{},
	}
}

// WithInput sets the input schema.
func (b *EndpointSchemaBuilder) WithInput(input *x402.InputSchema) *EndpointSchemaBuilder {
	b.schema.Input = input
	return b
}

// WithOutput sets the output schema.
func (b *EndpointSchemaBuilder) WithOutput(output map[string]any) *EndpointSchemaBuilder {
	b.schema.Output = output
	return b
}

// Build returns the constructed EndpointSchema.
func (b *EndpointSchemaBuilder) Build() *x402.EndpointSchema {
	return b.schema
}