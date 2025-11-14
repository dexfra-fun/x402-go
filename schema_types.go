package x402

import "encoding/json"

// FieldDef describes a field in the API schema according to x402 specification.
// It supports nested objects, enums, and conditional requirements.
type FieldDef struct {
	// Type specifies the field type (e.g., "string", "integer", "boolean", "object", "array").
	Type string `json:"type,omitempty"`

	// Required indicates if the field is required.
	// Can be a boolean (true/false) or an array of strings for conditional requirements.
	Required any `json:"required,omitempty"`

	// Description provides a human-readable description of the field.
	Description string `json:"description,omitempty"`

	// Enum lists allowed values for the field.
	Enum []string `json:"enum,omitempty"`

	// Properties defines nested fields for object types.
	Properties map[string]*FieldDef `json:"properties,omitempty"`
}

// InputSchema describes the input expectations for an API endpoint.
type InputSchema struct {
	// Type specifies the schema type (typically "http").
	Type string `json:"type"`

	// Method specifies the HTTP method (e.g., "GET", "POST", "PUT", "DELETE").
	Method string `json:"method"`

	// BodyType specifies the content type for request body.
	// Valid values: "json", "form-data", "multipart-form-data", "text", "binary".
	BodyType string `json:"bodyType,omitempty"`

	// QueryParams defines the query parameters accepted by the endpoint.
	QueryParams map[string]*FieldDef `json:"queryParams,omitempty"`

	// BodyFields defines the fields expected in the request body.
	BodyFields map[string]*FieldDef `json:"bodyFields,omitempty"`

	// HeaderFields defines the headers expected in the request.
	HeaderFields map[string]*FieldDef `json:"headerFields,omitempty"`

	// Discoverable indicates if the endpoint should be discoverable in schema listings.
	// Defaults to true if not specified (nil).
	Discoverable *bool `json:"discoverable,omitempty"`
}

// IsDiscoverable returns whether the endpoint is discoverable.
// Returns true by default if Discoverable is not explicitly set to false.
func (is *InputSchema) IsDiscoverable() bool {
	if is.Discoverable == nil {
		return true
	}
	return *is.Discoverable
}

// MarshalJSON implements custom JSON marshaling for InputSchema.
// It ensures Discoverable is set to true by default when not explicitly set.
func (is *InputSchema) MarshalJSON() ([]byte, error) {
	type Alias InputSchema
	
	// Create a copy with Discoverable set to default value if nil
	aux := &struct {
		Discoverable bool `json:"discoverable"`
		*Alias
	}{
		Discoverable: is.IsDiscoverable(),
		Alias:        (*Alias)(is),
	}
	
	return json.Marshal(aux)
}

// EndpointSchema wraps input and output schema for an API endpoint.
type EndpointSchema struct {
	// Input describes the expected input for the endpoint.
	Input *InputSchema `json:"input,omitempty"`

	// Output describes the expected output/response structure (optional).
	// This is a flexible map to support various output schemas.
	Output map[string]any `json:"output,omitempty"`
}