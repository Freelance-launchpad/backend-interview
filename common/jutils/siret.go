package jutils

import (
	"fmt"
	"regexp"
)

type Siret struct {
	Siren string
	NIC   string
}

var reSiret = regexp.MustCompile(`^\d{14}$`)

func NewSiret(siret string) (Siret, error) {
	if !reSiret.MatchString(siret) {
		return Siret{}, fmt.Errorf("invalid siret: %s", siret)
	}

	return Siret{
		Siren: siret[:9],
		NIC:   siret[9:],
	}, nil
}

func (s Siret) String() string {
	return s.Siren + s.NIC
}
