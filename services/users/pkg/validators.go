package pkg

import (
	"regexp"

	"github.com/teamwork/vat/v3"
)

var siretRegex = regexp.MustCompile(`^\d{14}$`)

// isSiretValid return true if the siret is only composed of 14 digits characters.
func isSiretValid(siret string) bool {
	return siretRegex.MatchString(siret)
}

// IsValid returns true if the legal information is valid.
func (li *LegalInformation) IsValid() bool {
	if li == nil || !li.LegalForm.IsValid() {
		return false
	}

	if li.SIRET == nil {
		return false
	}

	if !isSiretValid(*li.SIRET) {
		return false
	}

	if li.LegalForm == LegalFormEURL || li.LegalForm == LegalFormSASU {
		if li.CompanyName == nil || *li.CompanyName == "" {
			return false
		}
		if li.Capital == nil {
			return false
		}
	}

	if li.LiableVAT {
		if li.VATNumber == nil || vat.ValidateFormat(*li.VATNumber) != nil {
			return false
		}
		if li.NotLiableVATReason != nil {
			return false
		}
	} else {
		if li.NotLiableVATReason == nil || !li.NotLiableVATReason.IsValid() {
			return false
		}
	}

	return true
}

// IsValid returns true if the business information is valid.
func (bi *BusinessInformation) IsValid(allowDeprecatedProfessions bool) bool {
	if bi == nil || !bi.ProfessionType.IsValid(allowDeprecatedProfessions) {
		return false
	}

	switch bi.ProfessionType {
	case ProfessionTypeCraftsman:
		if bi.CCMACode == nil || *bi.CCMACode == "" {
			return false
		}
		return true
	default:
		if bi.CityOfRegistrationOffice == nil || *bi.CityOfRegistrationOffice == "" {
			return false
		}
		return true
	}
}

// IsValid returns true if the address is valid.
func (a *Address) IsValid() bool {
	if a == nil {
		return false
	}

	if a.Street == "" || a.City == "" || a.ZipCode == "" || a.Country == "" {
		return false
	}

	return true
}
