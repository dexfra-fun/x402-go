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

var (
	facilitatorIDToString = map[FacilitatorID]string{
		FacilitatorIDUnknown:         "",
		FacilitatorIDPayAI:           "payAI",
		FacilitatorIDCoinbase:        "coinbase",
		FacilitatorIDOpenX402:        "openx402",
		FacilitatorIDAurraCloud:      "aurracloud",
		FacilitatorIDCodeNut:         "codenut",
		FacilitatorIDCorbits:         "corbits",
		FacilitatorIDDaydreams:       "daydreams",
		FacilitatorIDDexter:          "dexter",
		FacilitatorIDUltravioletaDAO: "ultravioletadao",
	}

	stringToFacilitatorID = map[string]FacilitatorID{
		"payAI":           FacilitatorIDPayAI,
		"coinbase":        FacilitatorIDCoinbase,
		"openx402":        FacilitatorIDOpenX402,
		"aurracloud":      FacilitatorIDAurraCloud,
		"codenut":         FacilitatorIDCodeNut,
		"corbits":         FacilitatorIDCorbits,
		"daydreams":       FacilitatorIDDaydreams,
		"dexter":          FacilitatorIDDexter,
		"ultravioletadao": FacilitatorIDUltravioletaDAO,
	}
)

// IDString converts FacilitatorID to its string representation.
// This maps enum values to the actual facilitator IDs used in the registry.
func (f FacilitatorID) IDString() string {
	if str, ok := facilitatorIDToString[f]; ok {
		return str
	}
	return ""
}

// FacilitatorIDFromString converts a string ID to FacilitatorID enum.
func FacilitatorIDFromString(id string) FacilitatorID {
	if facilID, ok := stringToFacilitatorID[id]; ok {
		return facilID
	}
	return FacilitatorIDUnknown
}
