package facilitators

import x402 "github.com/dexfra-fun/x402-go"

// This file provides convenient access to all built-in facilitators.
// Import this package to automatically register all facilitators.

// All built-in facilitators are automatically registered via init() functions
// in their respective files (payai.go, openx402.go, etc.)

// GetByNetwork returns all facilitators that support the given network.
func GetByNetwork(network Network) []*Facilitator {
	return List(network)
}

// GetByID returns a specific facilitator by its string ID.
func GetByID(id string) *Facilitator {
	return Get(id)
}

// GetByEnum returns a specific facilitator by its FacilitatorID enum.
func GetByEnum(id x402.FacilitatorID) *Facilitator {
	return Get(id.IDString())
}

// GetAll returns all registered facilitators.
func GetAll() []*Facilitator {
	return List("")
}
