package pricing

import (
	"context"
	"testing"

	"github.com/dexfra-fun/x402-go/pkg/x402"
	"github.com/shopspring/decimal"
)

func TestFixed(t *testing.T) {
	tests := []struct {
		name     string
		amount   string
		expected string
	}{
		{"zero", "0", "0"},
		{"small", "0.001", "0.001"},
		{"large", "10.5", "10.5"},
		{"precise", "0.123456789", "0.123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount := decimal.RequireFromString(tt.amount)
			p := NewFixed(amount)
			got, err := p.GetPrice(context.Background(), x402.Resource{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got.String())
			}
		})
	}
}

func TestFixedFromFloat(t *testing.T) {
	p := NewFixedFromFloat(0.001)
	got, err := p.GetPrice(context.Background(), x402.Resource{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.String() != "0.001" {
		t.Errorf("expected 0.001, got %s", got.String())
	}
}

func TestPathBased(t *testing.T) {
	prices := map[string]float64{
		"/api/data":    0.001,
		"/api/premium": 0.01,
	}
	p := NewPathBasedFromFloat(prices, 0.005)

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"exact match data", "/api/data", "0.001"},
		{"exact match premium", "/api/premium", "0.01"},
		{"default", "/api/other", "0.005"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := x402.Resource{Path: tt.path}
			got, err := p.GetPrice(context.Background(), resource)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got.String())
			}
		})
	}
}

func TestMethodBased(t *testing.T) {
	prices := map[string]float64{
		"GET":  0.001,
		"POST": 0.005,
	}
	p := NewMethodBasedFromFloat(prices, 0.002)

	tests := []struct {
		name     string
		method   string
		expected string
	}{
		{"GET", "GET", "0.001"},
		{"POST", "POST", "0.005"},
		{"default PUT", "PUT", "0.002"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := x402.Resource{Method: tt.method}
			got, err := p.GetPrice(context.Background(), resource)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got.String())
			}
		})
	}
}
