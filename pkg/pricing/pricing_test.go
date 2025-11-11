package pricing

import (
	"context"
	"testing"

	"github.com/dexfra-fun/x402-go/pkg/x402"
)

func TestFixed(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		expected float64
	}{
		{"zero", 0, 0},
		{"small", 0.001, 0.001},
		{"large", 10.5, 10.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewFixed(tt.amount)
			got, err := p.GetPrice(context.Background(), x402.Resource{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, got)
			}
		})
	}
}

func TestPathBased(t *testing.T) {
	prices := map[string]float64{
		"/api/data":    0.001,
		"/api/premium": 0.01,
	}
	p := NewPathBased(prices, 0.005)

	tests := []struct {
		name     string
		path     string
		expected float64
	}{
		{"exact match data", "/api/data", 0.001},
		{"exact match premium", "/api/premium", 0.01},
		{"default", "/api/other", 0.005},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := x402.Resource{Path: tt.path}
			got, err := p.GetPrice(context.Background(), resource)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, got)
			}
		})
	}
}

func TestMethodBased(t *testing.T) {
	prices := map[string]float64{
		"GET":  0.001,
		"POST": 0.005,
	}
	p := NewMethodBased(prices, 0.002)

	tests := []struct {
		name     string
		method   string
		expected float64
	}{
		{"GET", "GET", 0.001},
		{"POST", "POST", 0.005},
		{"default PUT", "PUT", 0.002},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := x402.Resource{Method: tt.method}
			got, err := p.GetPrice(context.Background(), resource)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, got)
			}
		})
	}
}
