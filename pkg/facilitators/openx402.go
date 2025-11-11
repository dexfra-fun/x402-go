package facilitators

import "time"

var OpenX402 = &Facilitator{
	ID: "openx402",
	Metadata: FacilitatorMetadata{
		Name:    "OpenX402",
		DocsURL: "https://open.x402.host",
	},
	URL: "https://open.x402.host",
	Addresses: map[Network][]FacilitatorAddress{
		NetworkSolana: {
			{
				Address:                "5xvht4fYDs99yprfm4UeuHSLxMBRpotfBtUCQqM3oDNG",
				Tokens:                 []Token{USDCSolanaToken},
				DateOfFirstTransaction: time.Date(2025, 10, 16, 0, 0, 0, 0, time.UTC),
			},
		},
	},
}

func init() {
	Register(OpenX402)
}
