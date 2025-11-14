package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dexfra-fun/x402-go"
)

func main() {
	// Example 1: Discoverable not set - defaults to true
	schema1 := &x402.InputSchema{
		Type:   "http",
		Method: "GET",
	}

	// Example 2: Discoverable explicitly set to false
	falseVal := false
	schema2 := &x402.InputSchema{
		Type:         "http",
		Method:       "POST",
		Discoverable: &falseVal,
	}

	// Example 3: Discoverable explicitly set to true
	trueVal := true
	schema3 := &x402.InputSchema{
		Type:         "http",
		Method:       "PUT",
		Discoverable: &trueVal,
	}

	// Marshal and print
	if err := printSchema("Schema 1 (not set - defaults to true)", schema1); err != nil {
		log.Fatal(err)
	}
	if err := printSchema("Schema 2 (explicitly false)", schema2); err != nil {
		log.Fatal(err)
	}
	if err := printSchema("Schema 3 (explicitly true)", schema3); err != nil {
		log.Fatal(err)
	}
}

func printSchema(name string, schema *x402.InputSchema) error {
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", name, err)
	}
	fmt.Printf("\n%s:\n%s\n", name, string(data))
	fmt.Printf("IsDiscoverable() returns: %v\n", schema.IsDiscoverable())
	return nil
}