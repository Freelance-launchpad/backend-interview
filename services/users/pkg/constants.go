package pkg

import "slices"

type MaritalStatus string

const (
	MaritalStatusSingle      MaritalStatus = "single"
	MaritalStatusCivilUnion  MaritalStatus = "civil-union"
	MaritalStatusMaritalLife MaritalStatus = "marital-life"
	MaritalStatusMarried     MaritalStatus = "married"
	MaritalStatusDivorced    MaritalStatus = "divorced"
	MaritalStatusWidowed     MaritalStatus = "widowed"
	MaritalStatusSeparated   MaritalStatus = "separated"
)

func (ms MaritalStatus) IsValid() bool {
	switch ms {
	case MaritalStatusSingle, MaritalStatusCivilUnion, MaritalStatusMaritalLife, MaritalStatusMarried,
		MaritalStatusDivorced, MaritalStatusWidowed, MaritalStatusSeparated:
		return true
	default:
		return false
	}
}

// ProfessionType represents a profession type.
//
// Note: this enum should stay the same as Profession in marketing, except that merchant and craftsman
// are split in this one. This is because there are additional rules for the craftsman profession.
type ProfessionType string

const (
	ProfessionTypeBusinessConsultantFreelance  ProfessionType = "business_consultant_freelance"
	ProfessionTypeRealEstateProfessional       ProfessionType = "real_estate_professional"
	ProfessionTypeVTCDriver                    ProfessionType = "vtc_driver"
	ProfessionTypeProfessionInPersonalServices ProfessionType = "profession_in_personal_services"
	ProfessionTypeMerchant                     ProfessionType = "merchant"
	ProfessionTypeCraftsman                    ProfessionType = "craftsman"
	ProfessionTypeWellnessProfession           ProfessionType = "wellness_profession"
	ProfessionTypePhotographer                 ProfessionType = "photographer"
	// Métiers du spectacle et de la culture non intermittent
	ProfessionTypeEntertainmentProfession ProfessionType = "entertainment_profession"
	// Métiers réglementés (santé, juridique, architecture, fonction publique)
	ProfessionTypeRegulatedProfession ProfessionType = "regulated_profession"
	// Formateur indépendant
	ProfessionTypeFreelanceTrainer ProfessionType = "freelance_trainer"
	//nolint:misspell
	// Métiers de la Restauration, Hôtellerie, Tourisme
	ProfessionTypeHospitality   ProfessionType = "hospitality"
	ProfessionTypeConstruction  ProfessionType = "construction"
	ProfessionTypeHeavyIndustry ProfessionType = "heavy_industry"
	ProfessionTypeSecurity      ProfessionType = "security"

	// Deprecated: use ProfessionTypeFreelanceTrainer instead.
	ProfessionTypeQualiopiTrainer ProfessionType = "qualiopi_trainer"
	// Deprecated: use ProfessionTypeWellnessProfession or ProfessionTypeRegulatedProfession instead.
	ProfessionTypeParamedicalProfession ProfessionType = "paramedical_profession"
	// Deprecated: use ProfessionTypeBusinessConsultantFreelance instead.
	ProfessionTypeBusinessCoach ProfessionType = "business_coach"
)

// IsValid checks if the profession type is valid.
func (pt ProfessionType) IsValid(allowDeprecated bool) bool {
	if allowDeprecated && pt == ProfessionTypeParamedicalProfession {
		return true
	}

	return slices.Contains([]ProfessionType{
		ProfessionTypeBusinessConsultantFreelance,
		ProfessionTypeRealEstateProfessional,
		ProfessionTypeVTCDriver,
		ProfessionTypeProfessionInPersonalServices,
		ProfessionTypeMerchant,
		ProfessionTypeCraftsman,
		ProfessionTypeWellnessProfession,
		ProfessionTypePhotographer,
		ProfessionTypeEntertainmentProfession,
		ProfessionTypeRegulatedProfession,
		ProfessionTypeFreelanceTrainer,
		ProfessionTypeHospitality,
		ProfessionTypeConstruction,
		ProfessionTypeHeavyIndustry,
		ProfessionTypeSecurity,
	}, pt)
}

type LegalForm string

const (
	LegalFormMicroEnterprise    LegalForm = "micro_enterprise"
	LegalFormSoleProprietorship LegalForm = "sole_proprietorship"
	LegalFormEIRL               LegalForm = "eirl"
	LegalFormEURL               LegalForm = "eurl"
	LegalFormSASU               LegalForm = "sasu"
	LegalFormCAESCICSAS         LegalForm = "cae_scic_sas"
	// LegalFormLimitedCompany represents a UK limited company (Ltd).
	LegalFormLimitedCompany LegalForm = "limited_company"
)

func (lf LegalForm) IsValid() bool {
	return slices.Contains([]LegalForm{
		LegalFormMicroEnterprise,
		LegalFormSoleProprietorship,
		LegalFormEIRL,
		LegalFormEURL,
		LegalFormSASU,
		LegalFormCAESCICSAS,
		LegalFormLimitedCompany,
	}, lf)
}

func (lf LegalForm) String() string {
	switch lf {
	case LegalFormMicroEnterprise:
		return "Micro-entreprise"
	case LegalFormSoleProprietorship:
		return "Entreprise individuelle"
	case LegalFormEIRL:
		return "EIRL"
	case LegalFormEURL:
		return "EURL"
	case LegalFormSASU:
		return "SASU"
	case LegalFormCAESCICSAS:
		return "CAE SCIC SAS"
	case LegalFormLimitedCompany:
		return "Ltd"
	}
	return "--"
}

// AddressType represents an address type.
type AddressType string

const (
	AddressTypePersonal AddressType = "personal"
	AddressTypeBilling  AddressType = "billing"
)
