package facilitators

import "time"

var Daydreams = &Facilitator{
	ID: "daydreams",
	Metadata: FacilitatorMetadata{
		Name:    "Daydreams",
		DocsURL: "https://facilitator.daydreams.systems",
	},
	URL: "https://facilitator.daydreams.systems",
	Addresses: map[Network][]FacilitatorAddress{
		NetworkSolana: {
			{
				Address:                "DuQ4jFMmVABWGxabYHFkGzdyeJgS1hp4wrRuCtsJgT9a",
				Tokens:                 []Token{USDCSolanaToken},
				DateOfFirstTransaction: time.Date(2025, 10, 16, 0, 0, 0, 0, time.UTC),
			},
		},
	},
}

func init() {
	Register(Daydreams)
}
