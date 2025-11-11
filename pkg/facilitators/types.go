package facilitators

import "time"

// Network represents supported blockchain networks.
type Network string

const (
	NetworkSolana Network = "solana"
)

// Facilitator represents a complete facilitator configuration.
type Facilitator struct {
	ID        string
	Metadata  FacilitatorMetadata
	URL       string
	Addresses map[Network][]FacilitatorAddress
}

// FacilitatorMetadata contains display information for a facilitator.
type FacilitatorMetadata struct {
	Name    string
	DocsURL string
}

// FacilitatorAddress represents a facilitator's blockchain address with token support.
type FacilitatorAddress struct {
	Address                string
	Tokens                 []Token
	DateOfFirstTransaction time.Time
}

// Token represents a token configuration.
type Token struct {
	Address  string
	Decimals int
	Symbol   string
}

// Standard USDC token configuration for Solana
var USDCSolanaToken = Token{
	Address:  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
	Decimals: 6,
	Symbol:   "USDC",
}
