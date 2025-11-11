package main

import (
	"fmt"

	x402 "github.com/dexfra-fun/x402-go"
	"github.com/dexfra-fun/x402-go/pkg/facilitators"
)

func main() {
	fmt.Println("=== x402 Facilitator Enum Example ===")
	fmt.Println()

	// 1. Using enum directly
	fmt.Println("1. Get facilitator using enum:")
	id := x402.FacilitatorIDPayAI
	facilitator := facilitators.GetByEnum(id)
	if facilitator != nil {
		fmt.Printf("   - Enum: %s\n", id.String())
		fmt.Printf("   - ID: %s\n", id.IDString())
		fmt.Printf("   - Name: %s\n", facilitator.Metadata.Name)
		fmt.Printf("   - URL: %s\n", facilitator.URL)
	}

	// 2. Convert from string
	fmt.Println("\n2. Convert string to enum:")
	idFromString := x402.FacilitatorIDFromString("coinbase")
	if idFromString != x402.FacilitatorIDUnknown {
		fmt.Printf("   - Valid facilitator: %s\n", idFromString.IDString())
		facilitator := facilitators.GetByEnum(idFromString)
		fmt.Printf("   - Name: %s\n", facilitator.Metadata.Name)
	}

	// 3. List all facilitators using enum
	fmt.Println("\n3. All available facilitators:")
	allIDs := []x402.FacilitatorID{
		x402.FacilitatorIDPayAI,
		x402.FacilitatorIDCoinbase,
		x402.FacilitatorIDOpenX402,
		x402.FacilitatorIDAurraCloud,
		x402.FacilitatorIDCodeNut,
		x402.FacilitatorIDCorbits,
		x402.FacilitatorIDDaydreams,
		x402.FacilitatorIDDexter,
		x402.FacilitatorIDUltravioletaDAO,
	}

	for _, id := range allIDs {
		facilitator := facilitators.GetByEnum(id)
		if facilitator != nil {
			fmt.Printf("   - %s: %s\n", id.IDString(), facilitator.Metadata.Name)
		}
	}

	// 4. Safe enum usage
	fmt.Println("\n4. Safe enum handling:")
	unknownID := x402.FacilitatorIDFromString("invalid-id")
	if unknownID == x402.FacilitatorIDUnknown {
		fmt.Println("   - Safely handled unknown facilitator")
	}

	fmt.Println("\n✅ Example completed!")
}
