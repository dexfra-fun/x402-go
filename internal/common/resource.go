// Package common provides shared functionality for x402 middleware adapters.
package common

import (
	"net/http"
	"net/url"

	"github.com/dexfra-fun/x402-go/pkg/x402"
)

// ExtractResource creates a Resource from an HTTP request.
func ExtractResource(r *http.Request) x402.Resource {
	resource := x402.Resource{
		Path:   r.URL.Path,
		Method: r.Method,
		Params: make(map[string]string),
	}

	// Extract query parameters
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			resource.Params[key] = values[0]
		}
	}

	return resource
}

// ExtractResourceFromURL creates a Resource from URL and method.
func ExtractResourceFromURL(urlPath, method string, query url.Values) x402.Resource {
	resource := x402.Resource{
		Path:   urlPath,
		Method: method,
		Params: make(map[string]string),
	}

	// Extract query parameters
	for key, values := range query {
		if len(values) > 0 {
			resource.Params[key] = values[0]
		}
	}

	return resource
}
