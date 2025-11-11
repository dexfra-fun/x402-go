package facilitators

import "time"

// UltravioletaDAO is the Ultravioleta DAO facilitator configuration.
var UltravioletaDAO = &Facilitator{
	ID: "ultravioletadao",
	Metadata: FacilitatorMetadata{
		Name:    "Ultravioleta DAO",
		DocsURL: "https://facilitator.ultravioletadao.xyz",
	},
	URL: "https://facilitator.ultravioletadao.xyz",
	Addresses: map[Network][]FacilitatorAddress{
		NetworkSolana: {
			{
				Address:                "F742C4VfFLQ9zRQyithoj5229ZgtX2WqKCSFKgH2EThq",
				Tokens:                 []Token{USDCSolanaToken},
				DateOfFirstTransaction: time.Date(2025, 10, 30, 0, 0, 0, 0, time.UTC),
			},
		},
	},
}

func init() {
	if err := Register(UltravioletaDAO); err != nil {
		panic(err)
	}
}
