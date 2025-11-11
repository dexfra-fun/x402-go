package x402

import (
	"strings"

	"github.com/mark3labs/x402-go"
)

// MapNetworkToChain maps network string to x402.ChainConfig
func MapNetworkToChain(network string) (x402.ChainConfig, error) {
	normalized := strings.ToLower(strings.TrimSpace(network))
	
	switch normalized {
	case "solana-devnet":
		return x402.SolanaDevnet, nil
	case "solana-mainnet", "solana":
		return x402.SolanaMainnet, nil
	default:
		return x402.ChainConfig{}, ErrNetworkNotSupported
	}
}

// GetDefaultNetworks returns the default network configurations
func GetDefaultNetworks() map[string]NetworkConfig {
	return map[string]NetworkConfig{
		"solana-devnet": {
			ChainID:     "solana-devnet",
			Name:        "Solana Devnet",
			ChainConfig: x402.SolanaDevnet,
		},
		"solana-mainnet": {
			ChainID:     "solana-mainnet",
			Name:        "Solana Mainnet",
			ChainConfig: x402.SolanaMainnet,
		},
	}
}

// IsNetworkSupported checks if a network is supported
func IsNetworkSupported(network string) bool {
	normalized := strings.ToLower(strings.TrimSpace(network))
	networks := GetDefaultNetworks()
	
	for key := range networks {
		if key == normalized {
			return true
		}
	}
	
	return false
}
