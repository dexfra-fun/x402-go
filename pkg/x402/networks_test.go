package x402

import "testing"

func TestMapNetworkToChain(t *testing.T) {
	tests := []struct {
		name    string
		network string
		wantErr bool
	}{
		{"solana devnet", "solana-devnet", false},
		{"solana mainnet", "solana-mainnet", false},
		{"solana alias", "solana", false},
		{"unknown network", "unknown", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MapNetworkToChain(tt.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapNetworkToChain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsNetworkSupported(t *testing.T) {
	tests := []struct {
		name     string
		network  string
		expected bool
	}{
		{"solana devnet", "solana-devnet", true},
		{"solana mainnet", "solana-mainnet", true},
		{"unknown", "ethereum", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNetworkSupported(tt.network)
			if got != tt.expected {
				t.Errorf("IsNetworkSupported() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetDefaultNetworks(t *testing.T) {
	networks := GetDefaultNetworks()
	
	if len(networks) == 0 {
		t.Error("expected non-empty default networks")
	}
	
	if _, ok := networks["solana-devnet"]; !ok {
		t.Error("expected solana-devnet in default networks")
	}
	
	if _, ok := networks["solana-mainnet"]; !ok {
		t.Error("expected solana-mainnet in default networks")
	}
}
