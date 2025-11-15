package resource

import (
	"context"
	"testing"

	localx402 "github.com/dexfra-fun/x402-go/pkg/x402"
)

func TestPathBased_GetResourceURL(t *testing.T) {
	tests := []struct {
		name         string
		resources    map[string]*Metadata
		defaultMeta  *Metadata
		baseURL      string
		resourcePath string
		wantURL      string
	}{
		{
			name: "exact match with custom URL",
			resources: map[string]*Metadata{
				"/api/users": {
					URL:         "https://api.example.com/api/users",
					Description: "Users endpoint",
				},
			},
			resourcePath: "/api/users",
			wantURL:      "https://api.example.com/api/users",
		},
		{
			name:         "no match with base URL",
			resources:    map[string]*Metadata{},
			baseURL:      "https://api.example.com",
			resourcePath: "/api/products",
			wantURL:      "https://api.example.com/api/products",
		},
		{
			name: "no match with default resource",
			defaultMeta: &Metadata{
				URL: "https://default.example.com",
			},
			resourcePath: "/api/unknown",
			wantURL:      "https://default.example.com",
		},
		{
			name:         "no match no default no base URL",
			resources:    map[string]*Metadata{},
			resourcePath: "/api/test",
			wantURL:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPathBased(tt.resources, tt.defaultMeta, tt.baseURL)
			resource := localx402.Resource{Path: tt.resourcePath}

			gotURL, err := p.GetResourceURL(context.Background(), resource)
			if err != nil {
				t.Errorf("GetResourceURL() error = %v", err)
				return
			}
			if gotURL != tt.wantURL {
				t.Errorf("GetResourceURL() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func TestPathBased_GetDescription(t *testing.T) {
	tests := []struct {
		name            string
		resources       map[string]*Metadata
		defaultMeta     *Metadata
		resourcePath    string
		wantDescription string
	}{
		{
			name: "exact match with description",
			resources: map[string]*Metadata{
				"/api/users": {
					URL:         "https://api.example.com/api/users",
					Description: "Get list of users",
				},
			},
			resourcePath:    "/api/users",
			wantDescription: "Get list of users",
		},
		{
			name: "no match with default description",
			defaultMeta: &Metadata{
				Description: "Default API endpoint",
			},
			resourcePath:    "/api/unknown",
			wantDescription: "Default API endpoint",
		},
		{
			name:            "no match no default",
			resources:       map[string]*Metadata{},
			resourcePath:    "/api/test",
			wantDescription: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPathBased(tt.resources, tt.defaultMeta, "")
			resource := localx402.Resource{Path: tt.resourcePath}

			gotDesc, err := p.GetDescription(context.Background(), resource)
			if err != nil {
				t.Errorf("GetDescription() error = %v", err)
				return
			}
			if gotDesc != tt.wantDescription {
				t.Errorf("GetDescription() = %v, want %v", gotDesc, tt.wantDescription)
			}
		})
	}
}

func TestPathBased_AddResource(t *testing.T) {
	p := NewPathBased(nil, nil, "")

	// Add a resource
	p.AddResource("/api/test", &Metadata{
		URL:         "https://api.example.com/test",
		Description: "Test endpoint",
	})

	resource := localx402.Resource{Path: "/api/test"}

	// Check URL
	gotURL, err := p.GetResourceURL(context.Background(), resource)
	if err != nil {
		t.Errorf("GetResourceURL() error = %v", err)
		return
	}
	if gotURL != "https://api.example.com/test" {
		t.Errorf("GetResourceURL() = %v, want %v", gotURL, "https://api.example.com/test")
	}

	// Check description
	gotDesc, err := p.GetDescription(context.Background(), resource)
	if err != nil {
		t.Errorf("GetDescription() error = %v", err)
		return
	}
	if gotDesc != "Test endpoint" {
		t.Errorf("GetDescription() = %v, want %v", gotDesc, "Test endpoint")
	}
}

func TestPathBased_SetDefaultResource(t *testing.T) {
	p := NewPathBased(nil, nil, "")

	// Set default resource
	p.SetDefaultResource(&Metadata{
		URL:         "https://default.example.com",
		Description: "Default endpoint",
	})

	resource := localx402.Resource{Path: "/api/unknown"}

	// Check URL
	gotURL, err := p.GetResourceURL(context.Background(), resource)
	if err != nil {
		t.Errorf("GetResourceURL() error = %v", err)
		return
	}
	if gotURL != "https://default.example.com" {
		t.Errorf("GetResourceURL() = %v, want %v", gotURL, "https://default.example.com")
	}

	// Check description
	gotDesc, err := p.GetDescription(context.Background(), resource)
	if err != nil {
		t.Errorf("GetDescription() error = %v", err)
		return
	}
	if gotDesc != "Default endpoint" {
		t.Errorf("GetDescription() = %v, want %v", gotDesc, "Default endpoint")
	}
}

func TestPathBased_SetBaseURL(t *testing.T) {
	p := NewPathBased(nil, nil, "")

	// Set base URL
	p.SetBaseURL("https://api.example.com")

	resource := localx402.Resource{Path: "/api/test"}

	// Check URL is constructed from base URL
	gotURL, err := p.GetResourceURL(context.Background(), resource)
	if err != nil {
		t.Errorf("GetResourceURL() error = %v", err)
		return
	}
	expected := "https://api.example.com/api/test"
	if gotURL != expected {
		t.Errorf("GetResourceURL() = %v, want %v", gotURL, expected)
	}
}
