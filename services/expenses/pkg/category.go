package pkg

import "github.com/Freelance-launchpad/backend-interview/common/jentity"

type ExpenseCategory int

const (
	ExpenseCategoryPersonalMeal                   ExpenseCategory = 1
	ExpenseCategoryPhoneSubscription              ExpenseCategory = 2
	ExpenseCategoryInternetSubscription           ExpenseCategory = 3
	ExpenseCategoryBusinessMeal                   ExpenseCategory = 4
	ExpenseCategorySoftwareWebSubscriptionHosting ExpenseCategory = 5
	ExpenseCategoryNTICGear                       ExpenseCategory = 6
	ExpenseCategoryTransportSubscription          ExpenseCategory = 7
	ExpenseCategoryHousingLocation                ExpenseCategory = 10
	ExpenseCategoryCoworking                      ExpenseCategory = 11
	ExpenseCategoryGasoline                       ExpenseCategory = 12
	ExpenseCategoryToll                           ExpenseCategory = 13
	ExpenseCategoryParking                        ExpenseCategory = 14
	ExpenseCategoryOneTimeTransport               ExpenseCategory = 15
	ExpenseCategoryTaxi                           ExpenseCategory = 16
	ExpenseCategoryProfessionalTraining           ExpenseCategory = 18
	ExpenseCategoryVehicleRental                  ExpenseCategory = 19
	ExpenseCategoryConferences                    ExpenseCategory = 21
	ExpenseCategoryInternationalMeal              ExpenseCategory = 24
	ExpenseCategoryRent                           ExpenseCategory = 25
	ExpenseCategoryExpensesPackage                ExpenseCategory = 26
	ExpenseCategoryCoPropertyCharge               ExpenseCategory = 27
	ExpenseCategoryPropertyTax                    ExpenseCategory = 28
	ExpenseCategoryHouseInsurance                 ExpenseCategory = 29
	ExpenseCategoryIK                             ExpenseCategory = 30
	ExpenseCategoryHousingRenters                 ExpenseCategory = 31
	ExpenseCategoryHousingTenants                 ExpenseCategory = 32
	ExpenseCategoryOtherTransport                 ExpenseCategory = 33
)

const (
	ExpenseCategoryUKTravel                 ExpenseCategory = 34
	ExpenseCategoryUKSubsistenceMeals       ExpenseCategory = 35
	ExpenseCategoryUKParking                ExpenseCategory = 36
	ExpenseCategoryUKAccommodation          ExpenseCategory = 37
	ExpenseCategoryUKMaterials              ExpenseCategory = 38
	ExpenseCategoryUKTechSupplies           ExpenseCategory = 39
	ExpenseCategoryUKPPEWorkWear            ExpenseCategory = 40
	ExpenseCategoryUKSoftwareSubscriptions  ExpenseCategory = 41
	ExpenseCategoryUKTrainings              ExpenseCategory = 42
	ExpenseCategoryUKProfessionalMembership ExpenseCategory = 43
	ExpenseCategoryUKOther                  ExpenseCategory = 44
)

func RegionLanguageKey(r jentity.Region) string {
	switch r {
	case jentity.RegionUK:
		return "ENG"
	default:
		return "FRA"
	}
}

// PayslipExpenses are expenses that are paid as part of the salary payslip.
var PayslipExpenses = []ExpenseCategory{ExpenseCategoryTransportSubscription, ExpenseCategoryOtherTransport}

var ExpenseCategoriesForbiddenOnRestDays = []ExpenseCategory{
	ExpenseCategoryPersonalMeal,
	ExpenseCategoryBusinessMeal,
	ExpenseCategoryInternationalMeal,
	ExpenseCategoryToll,
	ExpenseCategoryTaxi,
	ExpenseCategoryIK,
	ExpenseCategoryVehicleRental,
	ExpenseCategoryParking,
	ExpenseCategoryOneTimeTransport,
	ExpenseCategoryCoworking,
	ExpenseCategoryProfessionalTraining,
	ExpenseCategoryConferences,
	ExpenseCategoryHousingLocation,
}

type Category struct {
	ID                ExpenseCategory `json:"id"`
	Icon              string          `json:"icon"`
	Name              string          `json:"name"`
	Selectable        bool            `json:"selectable"`
	ReimbursementRate int             `json:"reimbursement_rate"`
	SubCategories     []SubCategory   `json:"sub_categories"`
}

type ExpenseSubCategory string

const (
	// Train, TER, Intercité...
	ExpenseSubCategoryTransportSubscriptionRegional ExpenseSubCategory = "transport-subscription-regional"
	// Metro, bus, tramway...
	ExpenseSubCategoryTransportSubscriptionUrban ExpenseSubCategory = "transport-subscription-urban"
	// Bike, scooter...
	ExpenseSubCategoryTransportSubscriptionPersonal ExpenseSubCategory = "transport-subscription-personal"

	ExpenseSubCategoryOneTimeTransportTrain ExpenseSubCategory = "one-time-transport-train"
	ExpenseSubCategoryOneTimeTransportPlane ExpenseSubCategory = "one-time-transport-plane"

	ExpenseSubCategoryOtherTransportMetro        ExpenseSubCategory = "other-transport-metro"
	ExpenseSubCategoryOtherTransportBus          ExpenseSubCategory = "other-transport-bus"
	ExpenseSubCategoryOtherTransportTramway      ExpenseSubCategory = "other-transport-tramway"
	ExpenseSubCategoryOtherTransportBoat         ExpenseSubCategory = "other-transport-boat"
	ExpenseSubCategoryOtherTransportDiscountCard ExpenseSubCategory = "other-transport-discount-card"
)

type SubCategory struct {
	ID         ExpenseSubCategory `json:"id"`
	CategoryID ExpenseCategory    `json:"category_id"`
	Icon       string             `json:"icon"`
	Label      string             `json:"label"`
}
