package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isSiretValid(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		siret string
		want  bool
	}{
		"valid":               {siret: "12345678901234", want: true},
		"invalid - too short": {siret: "123456789012", want: false},
		"invalid - non-digit": {siret: "12345abc901234", want: false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, isSiretValid(tt.siret))
		})
	}
}

func Test_LegalInformation_IsValid(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		li   *LegalInformation
		want bool
	}{
		"valid": {
			li: &LegalInformation{
				LegalForm:          LegalFormMicroEnterprise,
				SIRET:              new("12345678901234"),
				LiableVAT:          false,
				NotLiableVATReason: new(NotLiableVATReasonExemptThreshold),
			},
			want: true,
		},
		"invalid - SIRET": {
			li: &LegalInformation{
				LegalForm: LegalFormMicroEnterprise,
				SIRET:     new("123"),
				LiableVAT: true,
			},
			want: false,
		},
		"invalid - missing CompanyName for SASU": {
			li: &LegalInformation{
				LegalForm:   LegalFormSASU,
				SIRET:       new("12345678901234"),
				CompanyName: nil,
				LiableVAT:   true,
			},
			want: false,
		},
		"invalid - missing Capital for SASU": {
			li: &LegalInformation{
				LegalForm:   LegalFormSASU,
				SIRET:       new("12345678901234"),
				CompanyName: new("Company"),
				LiableVAT:   true,
			},
			want: false,
		},
		"valid - VATNumber for liableVAT": {
			li: &LegalInformation{
				LegalForm: LegalFormMicroEnterprise,
				SIRET:     new("12345678901234"),
				LiableVAT: true,
				VATNumber: new("FROH807386388"),
			},
			want: true,
		},
		"invalid - VATNumber for liableVAT": {
			li: &LegalInformation{
				LegalForm: LegalFormMicroEnterprise,
				SIRET:     new("12345678901234"),
				LiableVAT: true,
				VATNumber: new("1234asdfasdfasdfasdf5678901"),
			},
			want: false,
		},
		"invalid - missing not liable vat reason": {
			li: &LegalInformation{
				LegalForm:          LegalFormMicroEnterprise,
				SIRET:              new("12345678901234"),
				LiableVAT:          false,
				NotLiableVATReason: nil,
			},
			want: false,
		},
		"invalid - not liable vat reason with liable vat": {
			li: &LegalInformation{
				LegalForm:          LegalFormMicroEnterprise,
				SIRET:              new("12345678901234"),
				LiableVAT:          true,
				NotLiableVATReason: new(NotLiableVATReasonExemptThreshold),
				VATNumber:          new("FRXX345678901"),
			},
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.li.IsValid())
		})
	}
}

func Test_BusinessInformation_IsValid(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		bi   *BusinessInformation
		want bool
	}{
		"valid - Craftsman profession": {
			bi: &BusinessInformation{
				ProfessionType: ProfessionTypeCraftsman,
				CCMACode:       new("1234"),
			},
			want: true,
		},
		"invalid - Craftsman missing CCMACode": {
			bi: &BusinessInformation{
				ProfessionType: ProfessionTypeCraftsman,
				CCMACode:       nil,
			},
			want: false,
		},
		"valid - other profession": {
			bi: &BusinessInformation{
				ProfessionType:           ProfessionTypeRealEstateProfessional,
				CityOfRegistrationOffice: new("Paris"),
			},
			want: true,
		},
		"invalid - missing city": {
			bi: &BusinessInformation{
				ProfessionType:           ProfessionTypeRealEstateProfessional,
				CityOfRegistrationOffice: nil,
			},
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.bi.IsValid(false))
		})
	}
}

func Test_Address_IsValid(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		a    *Address
		want bool
	}{
		"valid": {
			a: &Address{
				Street:  "123 Main St",
				City:    "Paris",
				ZipCode: "75000",
				Country: "France",
			},
			want: true,
		},
		"invalid - missing street": {
			a: &Address{
				City:    "Paris",
				ZipCode: "75000",
				Country: "France",
			},
			want: false,
		},
		"invalid - missing city": {
			a: &Address{
				Street:  "123 Main St",
				ZipCode: "75000",
				Country: "France",
			},
			want: false,
		},
		"invalid - missing country": {
			a: &Address{
				Street:  "123 Main St",
				City:    "Paris",
				ZipCode: "75000",
			},
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.a.IsValid())
		})
	}
}
