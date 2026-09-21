package jentity

import (
	"fmt"
	"slices"

	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
)

var ErrInvalidRegion = jerror.NewValidationError("INVALID_REGION")

type Region string

const (
	RegionFR Region = "fr"
	RegionUK Region = "uk"
)

func (r Region) DefaultLanguage() Language {
	switch r {
	case RegionFR:
		return LanguageFRFR
	case RegionUK:
		return LanguageENGB
	default:
		return LanguageFRFR
	}
}

func (r Region) IsValid() bool {
	return slices.Contains([]Region{RegionFR, RegionUK}, r)
}

// LifeEntities returns the life entities belonging to the region.
func (r Region) LifeEntities() []Entity {
	switch r {
	case RegionFR:
		return LifeEntities
	case RegionUK:
		return UKLifeEntities
	default:
		return nil
	}
}

func (r *Region) Scan(src any) error {
	switch v := src.(type) {
	case string:
		*r = Region(v)
	case []byte:
		*r = Region(v)
	default:
		return fmt.Errorf("unsupported type for Region: %T", src)
	}
	return nil
}
