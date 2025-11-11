package facilitators

import "time"

var Coinbase = &Facilitator{
	ID: "coinbase",
	Metadata: FacilitatorMetadata{
		Name:    "Coinbase",
		DocsURL: "https://docs.cdp.coinbase.com/x402/welcome",
	},
	URL: "https://facilitator.coinbase.com",
	Addresses: map[Network][]FacilitatorAddress{
		NetworkSolana: {
			{
				Address:                "L54zkaPQFeTn1UsEqieEXBqWrPShiaZEPD7mS5WXfQg",
				Tokens:                 []Token{USDCSolanaToken},
				DateOfFirstTransaction: time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC),
			},
		},
	},
}

func init() {
	Register(Coinbase)
}
