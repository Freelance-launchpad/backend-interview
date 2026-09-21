package pkg

import (
	"errors"
	"slices"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/google/uuid"
)

type AdditionalInformation struct {
	// Category 4 (business meals)
	Guests []GuestAdditionalInformation `json:"guests,omitempty"`

	// Category 6 (ntic)
	NTICGear *NTICGearAdditionalInformation `json:"ntic_gear,omitempty"`

	// Category 7 (transport subscription)
	TransportSubscription *TransportSubscriptionAdditionalInformation `json:"transport_subscription,omitempty"`

	// Category 10 (housing location)
	Location *LocationAdditionalInformation `json:"location,omitempty"`

	// Category 15 (one-time transport)
	OneTimeTransport *OneTimeTransportAdditionalInformation `json:"one_time_transport,omitempty"`

	// Category 19 (vehicle rental)
	Vehicle *VehicleAdditionalInformation `json:"vehicle,omitempty"`

	// Category 25 (rent)
	Rent *RentAdditionalInformation `json:"rent,omitempty"`

	// Category 27 (co-property charge)
	CopropertyTax *CopropertyTaxAdditionalInformation `json:"co_property_tax,omitempty"`

	// Category 28 (property tax)
	PropertyTax *PropertyTaxAdditionalInformation `json:"property_tax,omitempty"`

	// Category 29 (house insurance)
	HouseInsurance *HouseInsuranceAdditionalInformation `json:"house_insurance,omitempty"`

	// Category 30 (IK)
	IK *IKAdditionalInformation `json:"ik,omitempty"`

	// Category 31 (housing renters)
	HousingRenter *HousingRenterAdditionalInformation `json:"housing_renter,omitempty"`

	// Category 32 (housing tenants)
	HousingTenant *HousingTenantAdditionalInformation `json:"housing_tenant,omitempty"`

	// Category 33 (other transport)
	OtherTransport *OtherTransportAdditionalInformation `json:"other_transport,omitempty"`

	// UK expenses (categories 34-44)
	UKExpense *UKExpenseAdditionalInformation `json:"uk_expense,omitempty"`
}

type GuestAdditionalInformation struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company"`
	Role      string `json:"role"`
}

type NTICGearAdditionalInformation struct {
	Device NTICGearDevice `json:"device"`
}

type NTICGearDevice string

const (
	NTICGearDevicePhone    NTICGearDevice = "phone"
	NTICGearDeviceTablet   NTICGearDevice = "tablet"
	NTICGearDeviceComputer NTICGearDevice = "computer"
	NTICGearDeviceOther    NTICGearDevice = "other"
)

func (d NTICGearDevice) IsValid() bool {
	return slices.Contains([]NTICGearDevice{
		NTICGearDevicePhone,
		NTICGearDeviceTablet,
		NTICGearDeviceComputer,
		NTICGearDeviceOther,
	}, d)
}

type TransportSubscriptionAdditionalInformation struct {
	SubCategory  ExpenseSubCategory          `json:"sub_category"`
	Period       TransportSubscriptionPeriod `json:"period"`
	ValidityDate jdate.Date                  `json:"validity_date"`
	// DeclaredAmount contains the amount declared by the user, and should match the expense attachment.
	// In the case of annual subscription, the expense amount is computed as DeclaredAmount / 12.
	DeclaredAmount int64 `json:"declared_amount"`
}

type TransportSubscriptionPeriod string

const (
	TransportSubscriptionPeriodMonthly TransportSubscriptionPeriod = "monthly"
	TransportSubscriptionPeriodYearly  TransportSubscriptionPeriod = "yearly"
)

func (p TransportSubscriptionPeriod) IsValid() bool {
	return slices.Contains([]TransportSubscriptionPeriod{
		TransportSubscriptionPeriodMonthly,
		TransportSubscriptionPeriodYearly,
	}, p)
}

type LocationAdditionalInformation struct {
	Address string `json:"address"`
	ZipCode string `json:"zip_code"`
	City    string `json:"city"`
	Country string `json:"country"` // alpha-2 country code
}

type OneTimeTransportAdditionalInformation struct {
	SubCategory ExpenseSubCategory `json:"sub_category"`
}

type VehicleAdditionalInformation struct {
	LicensePlate string `json:"license_plate"`
}

type RentAdditionalInformation struct {
	TotalArea  float64 `json:"total_area"`
	WorkArea   float64 `json:"work_area"`
	RentAmount float64 `json:"rent_amount"`
}

type CopropertyTaxAdditionalInformation struct {
	TotalArea float64 `json:"total_area"`
	WorkArea  float64 `json:"work_area"`
	Amount    float64 `json:"co_property_tax_amount"`
}

type PropertyTaxAdditionalInformation struct {
	TotalArea float64 `json:"total_area"`
	WorkArea  float64 `json:"work_area"`
	Amount    float64 `json:"property_tax_amount"`
}

type HouseInsuranceAdditionalInformation struct {
	TotalArea float64 `json:"total_area"`
	WorkArea  float64 `json:"work_area"`
	Amount    float64 `json:"house_insurance_amount"`
}

type IKAdditionalInformation struct {
	VehicleID   uuid.UUID     `json:"vehicle_id"`
	Origin      IKCoordinates `json:"origin"`
	Destination IKCoordinates `json:"destination"`
	Length      uint64        `json:"length"`
	RoundTrip   bool          `json:"round_trip"`
}

type IKCoordinates struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type HousingRenterAdditionalInformation struct {
	TotalArea             float64 `json:"total_area" validate:"required,gt=0"`
	WorkArea              float64 `json:"work_area" validate:"required,gt=0"`
	MonthlyRentAmount     int64   `json:"monthly_rent_amount" validate:"required,gt=0"`
	YearlyInsuranceAmount int64   `json:"yearly_insurance_amount" validate:"gte=0"`
}

func (h HousingRenterAdditionalInformation) MonthlyCostByM2() float64 {
	return (float64(h.MonthlyRentAmount) +
		(float64(h.YearlyInsuranceAmount) / 12)) /
		h.TotalArea
}

type HousingTenantAdditionalInformation struct {
	YearlyRentalAmount         int64 `json:"yearly_rental_amount" validate:"required,gt=0"`
	YearlyPropertyTaxesAmount  int64 `json:"yearly_property_taxes" validate:"required,gt=0"`
	QuarterlyCondominiumAmount int64 `json:"quarterly_condominium_amount" validate:"gte=0"`
	YearlyHomeInsurance        int64 `json:"yearly_home_insurance" validate:"gte=0"`
}

func (h HousingTenantAdditionalInformation) MonthlyRentalAmount() float64 {
	return float64(h.YearlyRentalAmount) / 12
}

func (h HousingTenantAdditionalInformation) MonthlyPropertyTaxes() float64 {
	return float64(h.YearlyPropertyTaxesAmount) / 12
}

func (h HousingTenantAdditionalInformation) MonthlyCondominiumAmount() float64 {
	return float64(h.QuarterlyCondominiumAmount) / 3
}

func (h HousingTenantAdditionalInformation) MonthlyHomeInsurance() float64 {
	return float64(h.YearlyHomeInsurance) / 12
}

type OtherTransportAdditionalInformation struct {
	SubCategory ExpenseSubCategory `json:"sub_category"`
}

type UKExpenseAdditionalInformation struct {
	ApprovedInAssignment bool   `json:"approved_in_assignment"`
	RelationToAssignment string `json:"relation_to_assignment"`
}

// ErrMissingRelationToAssignment is returned when the uk_expense block has an empty relation_to_assignment.
var ErrMissingRelationToAssignment = errors.New("MISSING_RELATION_TO_ASSIGNMENT")

func (u UKExpenseAdditionalInformation) Validate() error {
	if u.RelationToAssignment == "" {
		return ErrMissingRelationToAssignment
	}
	return nil
}
