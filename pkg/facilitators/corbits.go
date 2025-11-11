package facilitators

import "time"

// Corbits is the Corbits facilitator configuration.
var Corbits = &Facilitator{
	ID: "corbits",
	Metadata: FacilitatorMetadata{
		Name:    "Corbits",
		DocsURL: "https://corbits.dev",
	},
	URL: "https://facilitator.corbits.dev",
	Addresses: map[Network][]FacilitatorAddress{
		NetworkSolana: {
			{
				Address:                "AepWpq3GQwL8CeKMtZyKtKPa7W91Coygh3ropAJapVdU",
				Tokens:                 []Token{USDCSolanaToken},
				DateOfFirstTransaction: time.Date(2025, 9, 21, 0, 0, 0, 0, time.UTC),
			},
		},
	},
}

func init() {
	if err := Register(Corbits); err != nil {
		panic(err)
	}
}
