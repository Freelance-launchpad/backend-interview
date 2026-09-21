package pkg

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/Freelance-launchpad/backend-interview/common/jdb"
	"github.com/Freelance-launchpad/backend-interview/common/jentity"
	"github.com/Freelance-launchpad/backend-interview/common/jutils"
	"github.com/google/uuid"
)

type UserWithOffersAndContracts struct {
	User   User                `json:"user"`
	Offers []OfferWithContract `json:"offers"`
}

type UserMarketing struct {
	UTMs map[string]string `json:"utms,omitempty"`
}

func (m UserMarketing) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *UserMarketing) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &m)
}

type User struct {
	ID                   string           `json:"id" binding:"required"`
	Email                string           `json:"email" binding:"required"`
	FirstName            string           `json:"first_name" binding:"required"`
	LastName             string           `json:"last_name" binding:"required"`
	Regions              []jentity.Region `json:"regions"`
	Language             jentity.Language `json:"language"`
	Timezone             jentity.Timezone `json:"timezone"`
	Phone                *jutils.Phone    `json:"phone"`
	BirthDate            *jdate.Date      `json:"birth_date"`
	SocialSecurityNumber *string          `json:"social_security_number"`
	MaritalName          *string          `json:"marital_name"`
	MaritalStatus        *MaritalStatus   `json:"marital_status"`
	BirthAddress         *BirthAddress    `json:"birth_address"`
	Nationality          *string          `json:"nationality"`
	Marketing            UserMarketing    `json:"marketing"`
	CreatedAt            time.Time        `json:"created_at"`
	UpdatedAt            *time.Time       `json:"updated_at"`
	CRMContactID         *string          `json:"crm_contact_id"`
	IntercomID           *string          `json:"intercom_id"`
	DisabledWorker       bool             `json:"disabled_worker"`
	LastSelectedOffer    *uuid.UUID       `json:"last_selected_offer_id"`
	SalesOwnerEmail      *string          `json:"sales_owner_email"`
	CSMOwnerEmail        *string          `json:"csm_owner_email"`
	DisabledAt           *time.Time       `json:"disabled_at"`
}

type Offer struct {
	ID          uuid.UUID         `json:"id"`
	UserID      string            `json:"user_id"`
	Type        jentity.OfferType `json:"type"`
	LegalInfoID *uuid.UUID        `json:"legal_info_id"`
	// FIXME: created_at is exposed on the API but it shouldn't be. We can remove it as soon as the front starts using enabled_at instead.
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"-"`
	CRMDealID  *string    `json:"crm_deal_id"`
	EnabledAt  *time.Time `json:"enabled_at"`
	DisabledAt *time.Time `json:"disabled_at"`
}

type OfferInfo struct {
	Entity       jentity.Entity `json:"entity"`
	IsOnboarding bool           `json:"is_onboarding"`
	IsActive     bool           `json:"is_active"`
}

type OfferWithContract struct {
	Offer    Offer     `json:"offer"`
	Contract *Contract `json:"contract"`
}

// GetEntity combine the offer and the contract (if a contract is attached to this offer)
// to determine the entity of the offer.
func (o OfferWithContract) GetEntity() jentity.Entity {
	if o.Offer.Type == jentity.OfferTypeOpen {
		return jentity.EntityOpen
	}

	if slices.Contains([]jentity.OfferType{
		jentity.OfferTypeLife,
		jentity.OfferTypeUKLife,
	}, o.Offer.Type) && o.Contract != nil {
		return o.Contract.Entity
	}

	return jentity.EntityOnboarding
}

// life: active during onboarding, not active when CDI ends.
// open: active when stripe is paid, inactive when subscription canceled
func (o OfferWithContract) IsActive() bool {
	if slices.Contains([]jentity.OfferType{
		jentity.OfferTypeLife,
		jentity.OfferTypeUKLife,
	}, o.Offer.Type) {
		if o.Contract != nil {
			if o.Contract.EndDate != nil {
				return jdate.Today().Before(*o.Contract.EndDate) || jdate.Today().Equal(*o.Contract.EndDate)
			}
		}
		return true
	}
	return o.Offer.EnabledAt != nil && o.Offer.DisabledAt == nil
}

// OfferWithContractSortBy is the type of a "less function" that defines the ordering of its OfferWithContract arguments.
type OfferWithContractSortBy func(p1, p2 *OfferWithContract) bool

// offerWithContractSorter joins a OfferWithContractSortBy function and a slice of OfferWithContract to be sorted.
// It follows the sort interface.
type offerWithContractSorter struct {
	offerWithContracts []OfferWithContract
	by                 OfferWithContractSortBy
}

// Len is part of sort.Interface.
func (s *offerWithContractSorter) Len() int {
	return len(s.offerWithContracts)
}

// Swap is part of sort.Interface.
func (s *offerWithContractSorter) Swap(i, j int) {
	s.offerWithContracts[i], s.offerWithContracts[j] = s.offerWithContracts[j], s.offerWithContracts[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *offerWithContractSorter) Less(i, j int) bool {
	return s.by(&s.offerWithContracts[i], &s.offerWithContracts[j])
}

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by OfferWithContractSortBy) Sort(offerWithContracts []OfferWithContract) {
	ps := &offerWithContractSorter{
		offerWithContracts: offerWithContracts,
		by:                 by,
	}
	sort.Sort(ps)
}

func rankOffer(o OfferWithContract) int {
	if slices.Contains([]jentity.OfferType{
		jentity.OfferTypeLife,
		jentity.OfferTypeUKLife,
	}, o.Offer.Type) {
		if o.IsActive() {
			// Live life - after onboarding
			if o.Contract != nil {
				return 1
			}
			// live life - onboarding
			return 2
		}
		// Ended life
		return 4
	}
	if o.Offer.Type == jentity.OfferTypeOpen {
		// Live open
		if o.IsActive() {
			return 3
		}
		// Ended open
		return 5
	}

	return 999
}

// FrontendDisplay sort the offer with contract following this order
// Live Life
// Onboarding Life
// Live Open
// Ended Life => se baser sur la fin du contrat
// Hint: smallest are first.
func FrontendDisplay(p1, p2 *OfferWithContract) bool {
	return rankOffer(*p1) < rankOffer(*p2)
}

type OfferMarketing struct {
	Origin OfferOrigin `json:"origin,omitempty"`
}

func (o OfferMarketing) Value() (driver.Value, error) {
	return json.Marshal(o)
}

type OfferOrigin string

const (
	// Unique Point d'Entrée
	OfferOriginUPE       OfferOrigin = "upe"
	OfferOriginFastTrack OfferOrigin = "fast_track"
	OfferOriginAdmin     OfferOrigin = "admin"
	// OpenOfferInfoPageApp is the page displayed to Life users to show them Open Offer information
	OfferOriginOpenOfferInfoPageApp OfferOrigin = "open_offer_info_page_app"
	// OpenOfferInfoPageApp is the page displayed to users without offer to show them Open Offer information
	OfferOriginOpenOfferInfoPageWeb OfferOrigin = "open_offer_info_page_web"
)

func (o OfferOrigin) IsValid() bool {
	switch o {
	case OfferOriginUPE, OfferOriginFastTrack, OfferOriginAdmin, OfferOriginOpenOfferInfoPageApp, OfferOriginOpenOfferInfoPageWeb:
		return true
	default:
		return false
	}
}

type LegalInfo struct {
	ID uuid.UUID `json:"id"`

	Legal          *LegalInformation    `json:"legal"`
	Business       *BusinessInformation `json:"business"`
	BillingAddress *Address             `json:"billing_address"`
}

type HRInfo struct {
	HRID    string    `json:"hr_id"`
	OfferID uuid.UUID `json:"offer_id"`
	UserID  string    `json:"user_id"`

	Entity    jentity.Entity `json:"entity"`
	BirthDate *jdate.Date    `json:"birth_date"`

	StartDate jdate.Date  `json:"start_date"`
	EndDate   *jdate.Date `json:"end_date"`
}

type GetHRInfoOutput struct {
	Count  uint64   `json:"count"`
	HRInfo []HRInfo `json:"hr_info"`
}

type GetHRIDsOuput struct {
	HRIDs []string `json:"hr_ids"`
}

type AccessModuleResultContext string

const (
	AccessModuleResultContextMissingFields AccessModuleResultContext = "MISSING_FIELDS"
)

type AccessModuleResult struct {
	Result        bool                       `json:"result"`
	Context       *AccessModuleResultContext `json:"context"`
	MissingFields []string                   `json:"missing_fields"`
}

type UserPersonalInfo struct {
	FirstName *string       `json:"first_name"`
	LastName  *string       `json:"last_name"`
	Phone     *jutils.Phone `json:"phone"`
}

// LegalInformation represents the legal information of an open user.
type LegalInformation struct {
	LegalForm   LegalForm `json:"legal_form" binding:"required"`
	CompanyName *string   `json:"company_name"`
	// CompanyRegistrationNumber is the identifying number of a company, eg.:
	// - SIREN in France (not SIRET, as it identifies an establishment of a company, and not a company)
	// - Company number in the UK
	CompanyRegistrationNumber *string             `json:"company_registration_number"`
	SIRET                     *string             `json:"siret"`
	Capital                   *int64              `json:"capital"`
	LiableVAT                 bool                `json:"liable_vat"`
	NotLiableVATReason        *NotLiableVATReason `json:"not_liable_vat_reason"`
	VATNumber                 *string             `json:"vat_number"`
}

type NotLiableVATReason string

const (
	// Franchise en base de TVA
	NotLiableVATReasonExemptThreshold NotLiableVATReason = "exempt_threshold"

	//nolint:misspell
	// Activité exonérée de TVA
	NotLiableVATReasonExemptActivity NotLiableVATReason = "exempt_activity"
)

func (r NotLiableVATReason) IsValid() bool {
	switch r {
	case NotLiableVATReasonExemptThreshold, NotLiableVATReasonExemptActivity:
		return true
	default:
		return false
	}
}

// Normalize removes all non-digit/letter characters from the Siret and VATNumber fields.
func (li *LegalInformation) Normalize() {
	if li == nil {
		return
	}

	if li.SIRET != nil && *li.SIRET != "" {
		var sb strings.Builder
		for _, c := range *li.SIRET {
			if unicode.IsDigit(c) {
				sb.WriteRune(c)
			}
		}
		li.SIRET = new(sb.String())
	}

	if li.VATNumber != nil {
		var sb strings.Builder
		for _, c := range *li.VATNumber {
			if unicode.IsDigit(c) || unicode.IsLetter(c) {
				sb.WriteRune(c)
			}
		}
		li.VATNumber = new(strings.ToUpper(sb.String()))
	}
}

// BusinessInformation represents the business information of an open user.
type BusinessInformation struct {
	ProfessionType           ProfessionType `json:"profession_type" binding:"required"`
	RCSNumber                *string        `json:"rcs_number"`
	RMNumber                 *string        `json:"rm_number"`
	CityOfRegistrationOffice *string        `json:"city_of_registration_office"`
	CCMACode                 *string        `json:"ccma_code"`
}

type BirthAddress struct {
	City    string `json:"city"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

type Address struct {
	ID         uuid.UUID   `json:"-"`
	Street     string      `json:"street" binding:"required"`
	Additional *string     `json:"additional"`
	City       string      `json:"city" binding:"required"`
	ZipCode    string      `json:"zip_code" binding:"required"`
	Country    string      `json:"country" binding:"required"`
	Type       AddressType `json:"type,omitempty"`
}

// String returns an address as a formatted string following this template: 'street - zip_code city - country'.
func (a Address) String() string {
	return fmt.Sprintf("%s - %s %s - %s", a.Street, a.ZipCode, a.City, a.Country)
}

type ProbationRenewal struct {
	EndDate      jdate.Date `json:"end_date"`
	RemoteUserID string     `json:"remote_user_id"`
}

type ProbationPeriod struct {
	EndDate jdate.Date        `json:"end_date"`
	Renewal *ProbationRenewal `json:"renewal"`
}

func (pp ProbationPeriod) Value() (driver.Value, error) {
	return json.Marshal(pp)
}

func (pp *ProbationPeriod) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &pp)
}

var ProbationMonths = map[jentity.Entity]int{
	jentity.EntityBlue:      4,
	jentity.EntityGreen:     3,
	jentity.EntityRealty:    3,
	jentity.EntityFormation: 3,
}

type Contract struct {
	ID              uuid.UUID        `json:"id"`
	ExternalID      string           `json:"external_id"`
	Type            string           `json:"type"`
	SignedAt        time.Time        `json:"signed_at"`
	StartDate       jdate.Date       `json:"start_date"`
	EndDate         *jdate.Date      `json:"end_date"`
	Entity          jentity.Entity   `json:"entity"`
	HRID            *string          `json:"hr_id"`
	OfferID         uuid.UUID        `json:"offer_id"`
	AddressID       uuid.UUID        `json:"address_id"`
	ProfessionID    *uuid.UUID       `json:"profession_id"`
	ProbationPeriod *ProbationPeriod `json:"probation_period"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       *time.Time       `json:"updated_at"`
}

// IsActive returns true if the contract is active at day 'd'.
func (c Contract) IsActive(d jdate.Date) bool {
	return !c.StartDate.After(d) && (c.EndDate == nil || !c.EndDate.Before(d))
}

type AttachContractV3Input struct {
	Contract Contract `json:"contract"`
	Address  Address  `json:"address"`
}

type RequestUpdateUser struct {
	FirstName            string         `json:"first_name" binding:"required"`
	LastName             string         `json:"last_name" binding:"required"`
	MaritalStatus        *MaritalStatus `json:"marital_status"`
	MaritalName          *string        `json:"marital_name"`
	SocialSecurityNumber *string        `json:"social_security_number"`
	Nationality          *string        `json:"nationality"`
	BirthDate            *jdate.Date    `json:"birth_date"`
	BirthAddress         *BirthAddress  `json:"birth_address"`
	DisabledWorker       bool           `json:"disabled_worker"`
	Phone                *jutils.Phone  `json:"phone"`
}

type RequestUpdateUserPersonalInfo struct {
	FirstName            string         `json:"first_name" binding:"required"`
	LastName             string         `json:"last_name" binding:"required"`
	MaritalStatus        *MaritalStatus `json:"marital_status"`
	MaritalName          *string        `json:"marital_name"`
	SocialSecurityNumber *string        `json:"social_security_number"`
	Nationality          *string        `json:"nationality"`
	BirthDate            *jdate.Date    `json:"birth_date"`
	BirthAddress         *BirthAddress  `json:"birth_address"`
	DisabledWorker       bool           `json:"disabled_worker"`
}

type RequestUpdateIntercomID struct {
	IntercomID string `json:"intercom_id"`
}

type RequestUpdateUserEmail struct {
	Email string `json:"email"`
}

type UserInfo struct {
	ID              string                 `json:"id"`
	CreatedAt       time.Time              `json:"created_at"`
	FirstName       string                 `json:"first_name"`
	LastName        string                 `json:"last_name"`
	Email           string                 `json:"email"`
	Regions         []jentity.Region       `json:"regions"`
	Phone           *jutils.Phone          `json:"phone"`
	Offers          []OfferAndContractInfo `json:"offers"`
	SalesOwnerEmail *string                `json:"sales_owner_email"`
	CSMOwnerEmail   *string                `json:"csm_owner_email"`
	IntercomID      *string                `json:"intercom_id"`
}

// OfferAndContractInfo can contain offer info in the case of an open offer, or contract info in the case of a life offer.
type OfferAndContractInfo struct {
	ID              uuid.UUID        `json:"id"`
	Entity          jentity.Entity   `json:"entity"`
	StartDate       jdate.Date       `json:"start_date"`
	EndDate         *jdate.Date      `json:"end_date"`
	ProbationPeriod *ProbationPeriod `json:"probation_period"`
	DisabledAt      *time.Time       `json:"disabled_at"`
}

func (i OfferAndContractInfo) IsLiveAt(d jdate.Date) bool {
	if i.StartDate.After(d) {
		return false
	}

	if i.EndDate != nil && i.EndDate.Before(d) {
		return false
	}

	return true
}

type LifeUserInfo struct {
	UserID          string          `json:"user_id"`
	OfferID         uuid.UUID       `json:"offer_id"`
	HRID            *string         `json:"hr_id"`
	FirstName       string          `json:"first_name"`
	LastName        string          `json:"last_name"`
	Email           string          `json:"email"`
	BirthDate       *jdate.Date     `json:"birth_date"`
	Entity          *jentity.Entity `json:"entity"`
	StartDate       *jdate.Date     `json:"start_date"`
	EndDate         *jdate.Date     `json:"end_date"`
	SalesOwnerEmail *string         `json:"sales_owner_email"`
	CSMOwnerEmail   *string         `json:"csm_owner_email"`
}

type ListUsersInfoParams struct {
	jdb.QueryPagination
	Search          string                    `form:"search"`
	Offer           ListUsersInfoParamsOffer  `form:"offer"`
	Status          ListUsersInfoParamsStatus `form:"status"`
	Region          jentity.Region            `form:"region"`
	SalesOwnerEmail string                    `form:"sales_owner_email"`
	CSMOwnerEmail   string                    `form:"csm_owner_email"`
	HasCRMID        *bool                     `form:"has_crm_id"`
	IncludeArchived bool                      `form:"include_archived"`
}

func (p ListUsersInfoParams) AddToQuery(uri url.URL) url.URL {
	uri = p.AddURLValues(uri)

	values := uri.Query()

	if p.Search != "" {
		values.Add("search", p.Search)
	}
	if p.Offer != "" {
		values.Add("offer", string(p.Offer))
	}
	if p.Status != "" {
		values.Add("status", string(p.Status))
	}
	if p.Region != "" {
		values.Add("region", string(p.Region))
	}
	if p.SalesOwnerEmail != "" {
		values.Add("sales_owner_email", p.SalesOwnerEmail)
	}
	if p.CSMOwnerEmail != "" {
		values.Add("csm_owner_email", p.CSMOwnerEmail)
	}
	if p.HasCRMID != nil {
		values.Add("has_crm_id", strconv.FormatBool(*p.HasCRMID))
	}
	if p.IncludeArchived {
		values.Add("include_archived", "true")
	}

	uri.RawQuery = values.Encode()
	return uri
}

// ListUsersInfoParamsOffer represents the entity of a user. It's used instead of [jentity.Entity] because we don't want to expose deprecated entities such as sap or driver.
type ListUsersInfoParamsOffer string

const (
	ListUsersInfoParamsOfferBlue      ListUsersInfoParamsOffer = ListUsersInfoParamsOffer(jentity.EntityBlue)
	ListUsersInfoParamsOfferGreen     ListUsersInfoParamsOffer = ListUsersInfoParamsOffer(jentity.EntityGreen)
	ListUsersInfoParamsOfferRealty    ListUsersInfoParamsOffer = ListUsersInfoParamsOffer(jentity.EntityRealty)
	ListUsersInfoParamsOfferFormation ListUsersInfoParamsOffer = ListUsersInfoParamsOffer(jentity.EntityFormation)
	ListUsersInfoParamsOfferAtWork    ListUsersInfoParamsOffer = ListUsersInfoParamsOffer(jentity.EntityAtWork)
	ListUsersInfoParamsOfferOpen      ListUsersInfoParamsOffer = ListUsersInfoParamsOffer(jentity.EntityOpen)
	ListUsersInfoParamsOfferNone      ListUsersInfoParamsOffer = "none"
)

type ListUsersInfoParamsStatus string

const (
	ListUsersInfoParamsStatusLive       ListUsersInfoParamsStatus = "live"
	ListUsersInfoParamsStatusNotStarted ListUsersInfoParamsStatus = "not_started"
	ListUsersInfoParamsStatusEnded      ListUsersInfoParamsStatus = "ended"
)

type ResponseListUsersInfo struct {
	Count uint64     `json:"count"`
	Users []UserInfo `json:"users"`
}

type FileInfosMetadata struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type FileInfo struct {
	ID             uuid.UUID           `json:"id"`
	OfferID        uuid.UUID           `json:"offer_id"`
	Name           string              `json:"name"`
	RemoteID       string              `json:"remote_id"`
	Tag            *string             `json:"tag"`
	URL            string              `json:"url"`
	AdminOnly      bool                `json:"-"`
	ExpirationDate *jdate.Date         `json:"expiration_date"`
	Metadatas      []FileInfosMetadata `json:"metadatas"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      *time.Time          `json:"updated_at"`
}
