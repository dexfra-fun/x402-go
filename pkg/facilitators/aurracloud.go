package facilitators

import "time"

var AurraCloud = &Facilitator{
	ID: "aurracloud",
	Metadata: FacilitatorMetadata{
		Name:    "AurraCloud",
		DocsURL: "https://x402-facilitator.aurracloud.com",
	},
	URL: "https://x402-facilitator.aurracloud.com",
	Addresses: map[Network][]FacilitatorAddress{
		NetworkSolana: {
			{
				Address:                "8x8CzkTHTYkW18frrTR7HdCV6fsjenvcykJAXWvoPQW",
				Tokens:                 []Token{USDCSolanaToken},
				DateOfFirstTransaction: time.Date(2025, 10, 30, 0, 0, 0, 0, time.UTC),
			},
		},
	},
}

func init() {
	Register(AurraCloud)
}
