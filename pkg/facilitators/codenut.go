package facilitators

import "time"

var CodeNut = &Facilitator{
	ID: "codenut",
	Metadata: FacilitatorMetadata{
		Name:    "CodeNut",
		DocsURL: "https://docs.codenut.ai/guides/x402-facilitator",
	},
	URL: "https://facilitator.codenut.ai",
	Addresses: map[Network][]FacilitatorAddress{
		NetworkSolana: {
			{
				Address:                "HsozMJWWHNADoZRmhDGKzua6XW6NNfNDdQ4CkE9i5wHt",
				Tokens:                 []Token{USDCSolanaToken},
				DateOfFirstTransaction: time.Date(2025, 11, 3, 0, 0, 0, 0, time.UTC),
			},
		},
	},
}

func init() {
	Register(CodeNut)
}
