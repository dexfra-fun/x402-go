package facilitators

import "time"

// PayAI is the PayAI facilitator configuration.
var PayAI = &Facilitator{
	ID: "payAI",
	Metadata: FacilitatorMetadata{
		Name:    "PayAI",
		DocsURL: "https://payai.network",
	},
	URL: "https://facilitator.payai.network",
	Addresses: map[Network][]FacilitatorAddress{
		NetworkSolana: {
			{
				Address:                "2wKupLR9q6wXYppw8Gr2NvWxKBUqm4PPJKkQfoxHDBg4",
				Tokens:                 []Token{USDCSolanaToken},
				DateOfFirstTransaction: time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	},
}

func init() {
	if err := Register(PayAI); err != nil {
		panic(err)
	}
}
