package jentity

import "slices"

type OfferType string

const (
	// FR

	OfferTypeLife OfferType = "life"
	// TODO: rename Open as Micro
	OfferTypeOpen OfferType = "open"

	// UK

	OfferTypeUKLife OfferType = "uk-life"
)

func (ot OfferType) HasOnboarding() bool {
	return slices.Contains([]OfferType{
		OfferTypeLife,
		OfferTypeUKLife,
	}, ot)
}
