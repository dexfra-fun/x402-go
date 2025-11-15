package resource

import (
	"context"
	"fmt"

	localx402 "github.com/dexfra-fun/x402-go/pkg/x402"
)

// PathBased provides different resource URLs and descriptions based on the request path.
// This allows different endpoints to have different resource metadata.
type PathBased struct {
	resources       map[string]*Metadata
	defaultResource *Metadata
	baseURL         string
}

// Metadata contains resource URL and description.
type Metadata struct {
	URL         string
	Description string
}

// NewPathBased creates a new path-based resource provider.
// The resources map keys should be the path patterns to match.
// If defaultResource is provided, it will be used when no path matches.
// If baseURL is provided, it will be used to construct full URLs from relative paths.
func NewPathBased(resources map[string]*Metadata, defaultResource *Metadata, baseURL string) *PathBased {
	return &PathBased{
		resources:       resources,
		defaultResource: defaultResource,
		baseURL:         baseURL,
	}
}

// GetResourceURL returns the resource URL for the given resource path.
// If no matching resource is found, returns the default resource URL (if configured).
// If baseURL is configured, it constructs the full URL from the path.
func (p *PathBased) GetResourceURL(_ context.Context, resource localx402.Resource) (string, error) {
	// Try exact path match first
	if metadata, ok := p.resources[resource.Path]; ok {
		if metadata.URL != "" {
			return metadata.URL, nil
		}
	}

	// Use default resource if configured
	if p.defaultResource != nil && p.defaultResource.URL != "" {
		return p.defaultResource.URL, nil
	}

	// Construct URL from baseURL and path
	if p.baseURL != "" {
		return fmt.Sprintf("%s%s", p.baseURL, resource.Path), nil
	}

	return "", nil
}

// GetDescription returns the description for the given resource path.
// If no matching resource is found, returns the default description (if configured).
func (p *PathBased) GetDescription(_ context.Context, resource localx402.Resource) (string, error) {
	// Try exact path match first
	if metadata, ok := p.resources[resource.Path]; ok {
		if metadata.Description != "" {
			return metadata.Description, nil
		}
	}

	// Use default description if configured
	if p.defaultResource != nil {
		return p.defaultResource.Description, nil
	}

	return "", nil
}

// AddResource adds or updates resource metadata for a specific path.
func (p *PathBased) AddResource(path string, metadata *Metadata) {
	if p.resources == nil {
		p.resources = make(map[string]*Metadata)
	}
	p.resources[path] = metadata
}

// SetDefaultResource sets the default resource metadata to use when no path matches.
func (p *PathBased) SetDefaultResource(metadata *Metadata) {
	p.defaultResource = metadata
}

// SetBaseURL sets the base URL for constructing full resource URLs.
func (p *PathBased) SetBaseURL(baseURL string) {
	p.baseURL = baseURL
}
