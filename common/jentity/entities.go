package jentity

import (
	"encoding/json"
	"slices"

	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
)

type Entity string

const (
	// FR

	EntityBlue       Entity = "blue"
	EntityGreen      Entity = "green"
	EntityRealty     Entity = "realty"
	EntityFormation  Entity = "formation"
	EntityOpen       Entity = "open"
	EntityOnboarding Entity = "onboarding"

	// UK

	EntityAtWork Entity = "at-work"

	// Deprecated
	EntitySAP Entity = "sap"
)

var (
	LifeEntities               = []Entity{EntityBlue, EntityGreen, EntityRealty, EntityFormation}
	UKLifeEntities             = []Entity{EntityAtWork}
	Entities                   = append(LifeEntities, EntityOpen)
	EntitiesWithOnboarding     = append(Entities, EntityOnboarding)
	LifeEntitiesWithOnboarding = append(LifeEntities, EntityOnboarding)

	ErrInvalidEntity = jerror.NewValidationError("INVALID_ENTITY")
)

func isValidEntity(entity string) error {
	if !slices.Contains([]Entity{
		EntityBlue,
		EntityGreen,
		EntityRealty,
		EntityFormation,
		EntityOpen,
		EntityOnboarding,
		EntityAtWork,
		EntitySAP,
	}, Entity(entity)) {
		return ErrInvalidEntity
	}

	return nil
}

func NewEntity(entity string) (Entity, error) {
	if err := isValidEntity(entity); err != nil {
		return Entity(""), err
	}
	return Entity(entity), nil
}

func (e *Entity) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	entity, err := NewEntity(s)
	if err != nil {
		return err
	}

	*e = entity
	return nil
}

func (e Entity) IsValid() bool {
	err := isValidEntity(string(e))
	return err == nil
}

func (e Entity) IsLife() bool {
	return slices.Contains(LifeEntities, e)
}

func (e Entity) IsUKLife() bool {
	return slices.Contains(UKLifeEntities, e)
}

func (e Entity) IsLifeOrOnboarding() bool {
	return slices.Contains(LifeEntitiesWithOnboarding, e)
}

func (e Entity) String() string {
	return string(e)
}

func (e Entity) HasOnboarding() bool {
	return e.IsLife() || e.IsUKLife()
}

func (e Entity) HasJobContractStartDate() bool {
	return slices.Contains([]Entity{
		EntityGreen,
		EntityRealty,
	}, e)
}

func (e Entity) Region() Region {
	switch {
	case slices.Contains(LifeEntities, e):
		return RegionFR
	case slices.Contains(UKLifeEntities, e):
		return RegionUK
	default:
		return ""
	}
}
