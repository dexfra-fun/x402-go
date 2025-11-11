// Package x402 provides core types and enums for the x402 payment protocol.
package x402

// FacilitatorID represents unique identifiers for x402 facilitators.
//
//go:generate enumer -type=FacilitatorID -json -text -yaml -sql -transform=snake
type FacilitatorID int

const (
	// FacilitatorIDUnknown represents an unknown facilitator.
	FacilitatorIDUnknown FacilitatorID = iota
	// FacilitatorIDPayAI represents PayAI facilitator.
	FacilitatorIDPayAI
	// FacilitatorIDCoinbase represents Coinbase facilitator.
	FacilitatorIDCoinbase
	// FacilitatorIDOpenX402 represents OpenX402 facilitator.
	FacilitatorIDOpenX402
	// FacilitatorIDAurraCloud represents AurraCloud facilitator.
	FacilitatorIDAurraCloud
	// FacilitatorIDCodeNut represents CodeNut facilitator.
	FacilitatorIDCodeNut
	// FacilitatorIDCorbits represents Corbits facilitator.
	FacilitatorIDCorbits
	// FacilitatorIDDaydreams represents Daydreams facilitator.
	FacilitatorIDDaydreams
	// FacilitatorIDDexter represents Dexter facilitator.
	FacilitatorIDDexter
	// FacilitatorIDUltravioletaDAO represents UltravioletaDAO facilitator.
	FacilitatorIDUltravioletaDAO
)

// IDString converts FacilitatorID to its string representation.
// This maps enum values to the actual facilitator IDs used in the registry.
func (f FacilitatorID) IDString() string {
	switch f {
	case FacilitatorIDPayAI:
		return "payAI"
	case FacilitatorIDCoinbase:
		return "coinbase"
	case FacilitatorIDOpenX402:
		return "openx402"
	case FacilitatorIDAurraCloud:
		return "aurracloud"
	case FacilitatorIDCodeNut:
		return "codenut"
	case FacilitatorIDCorbits:
		return "corbits"
	case FacilitatorIDDaydreams:
		return "daydreams"
	case FacilitatorIDDexter:
		return "dexter"
	case FacilitatorIDUltravioletaDAO:
		return "ultravioletadao"
	default:
		return ""
	}
}

// FacilitatorIDFromString converts a string ID to FacilitatorID enum.
func FacilitatorIDFromString(id string) FacilitatorID {
	switch id {
	case "payAI":
		return FacilitatorIDPayAI
	case "coinbase":
		return FacilitatorIDCoinbase
	case "openx402":
		return FacilitatorIDOpenX402
	case "aurracloud":
		return FacilitatorIDAurraCloud
	case "codenut":
		return FacilitatorIDCodeNut
	case "corbits":
		return FacilitatorIDCorbits
	case "daydreams":
		return FacilitatorIDDaydreams
	case "dexter":
		return FacilitatorIDDexter
	case "ultravioletadao":
		return FacilitatorIDUltravioletaDAO
	default:
		return FacilitatorIDUnknown
	}
}
