package facilitators

import "time"

// Dexter is the Dexter facilitator configuration.
var Dexter = &Facilitator{
	ID: "dexter",
	Metadata: FacilitatorMetadata{
		Name:    "Dexter",
		DocsURL: "https://facilitator.dexter.cash",
	},
	URL: "https://facilitator.dexter.cash",
	Addresses: map[Network][]FacilitatorAddress{
		NetworkSolana: {
			{
				Address:                "DEXVS3su4dZQWTvvPnLDJLRK1CeeKG6K3QqdzthgAkNV",
				Tokens:                 []Token{USDCSolanaToken},
				DateOfFirstTransaction: time.Date(2025, 10, 26, 0, 0, 0, 0, time.UTC),
			},
		},
	},
}

func init() {
	if err := Register(Dexter); err != nil {
		panic(err)
	}
}
