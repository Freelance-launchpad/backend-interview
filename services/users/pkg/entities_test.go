package pkg

import (
	"net/url"
	"testing"
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/Freelance-launchpad/backend-interview/common/jdb"
	"github.com/Freelance-launchpad/backend-interview/common/jentity"
	"github.com/stretchr/testify/assert"
)

func Test_OfferWithContract_GetEntity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		obj    OfferWithContract
		output jentity.Entity
	}{
		{
			name: "open",
			obj: OfferWithContract{
				Offer: Offer{Type: jentity.OfferTypeOpen},
			},
			output: jentity.EntityOpen,
		},
		{
			name: "green",
			obj: OfferWithContract{
				Offer:    Offer{Type: jentity.OfferTypeLife},
				Contract: &Contract{Entity: jentity.EntityGreen},
			},
			output: jentity.EntityGreen,
		},
		{
			name: "blue",
			obj: OfferWithContract{
				Offer:    Offer{Type: jentity.OfferTypeLife},
				Contract: &Contract{Entity: jentity.EntityBlue},
			},
			output: jentity.EntityBlue,
		},
		{
			name: "unknown",
			obj: OfferWithContract{
				Offer: Offer{Type: jentity.OfferTypeLife},
			},
			output: jentity.EntityOnboarding,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.output, tt.obj.GetEntity())
		})
	}
}

func Test_OfferWithContract_IsActive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		offerWithContract OfferWithContract
		output            bool
	}{
		{
			name: "life active",
			offerWithContract: OfferWithContract{
				Offer:    Offer{Type: jentity.OfferTypeLife},
				Contract: &Contract{Type: jentity.EntityGreen.String()},
			},
			output: true,
		},
		{
			name: "life active last day",
			offerWithContract: OfferWithContract{
				Offer:    Offer{Type: jentity.OfferTypeLife},
				Contract: &Contract{Type: jentity.EntityGreen.String(), EndDate: new(jdate.Today())},
			},
			output: true,
		},
		{
			name: "life active - onboarding",
			offerWithContract: OfferWithContract{
				Offer: Offer{Type: jentity.OfferTypeLife},
			},
			output: true,
		},
		{
			name: "open active",
			offerWithContract: OfferWithContract{
				Offer: Offer{
					EnabledAt: new(time.Now()),
				},
			},
			output: true,
		},
		{
			name: "open not started",
			offerWithContract: OfferWithContract{
				Offer: Offer{},
			},
		},
		{
			name: "open ended",
			offerWithContract: OfferWithContract{
				Offer: Offer{
					DisabledAt: new(time.Now()),
				},
			},
			output: false,
		},
		{
			name: "life ended - contract ended",
			offerWithContract: OfferWithContract{
				Offer:    Offer{Type: jentity.OfferTypeLife},
				Contract: &Contract{Type: jentity.EntityGreen.String(), EndDate: new(jdate.New(2023, 1, 1))},
			},
			output: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.output, tt.offerWithContract.IsActive())
		})
	}
}

func Test_OfferWithContractSortBy_Sort(t *testing.T) {
	t.Parallel()

	disabledAt := time.Date(2024, time.February, 3, 0, 0, 0, 0, time.UTC)
	enabledAt := disabledAt.Add(time.Duration(-24 * time.Hour))

	tests := []struct {
		name     string
		original []OfferWithContract
		result   []OfferWithContract
	}{
		{
			name: "full potential",
			original: []OfferWithContract{
				/* Open live             */ {Offer: Offer{Type: jentity.OfferTypeOpen, EnabledAt: &enabledAt}},
				/* Open Ended            */ {Offer: Offer{Type: jentity.OfferTypeOpen, EnabledAt: &enabledAt, DisabledAt: &disabledAt}},
				/* Life live futur churn */ {Offer: Offer{Type: jentity.OfferTypeLife}, Contract: &Contract{EndDate: new(jdate.NewFromTime(time.Now().Add(48 * time.Hour)))}},
				/* Life Ended            */ {Offer: Offer{Type: jentity.OfferTypeLife}, Contract: &Contract{EndDate: new(jdate.NewFromTime(time.Now().Add(-48 * time.Hour)))}},
				/* life onboarding       */ {Offer: Offer{Type: jentity.OfferTypeLife}},
				/* life live             */ {Offer: Offer{Type: jentity.OfferTypeLife}, Contract: &Contract{}},
			},
			result: []OfferWithContract{
				/* Life live futur churn */ {Offer: Offer{Type: jentity.OfferTypeLife}, Contract: &Contract{EndDate: new(jdate.NewFromTime(time.Now().Add(48 * time.Hour)))}},
				/* life live             */ {Offer: Offer{Type: jentity.OfferTypeLife}, Contract: &Contract{}},
				/* life onboarding       */ {Offer: Offer{Type: jentity.OfferTypeLife}},
				/* Open live             */ {Offer: Offer{Type: jentity.OfferTypeOpen, EnabledAt: &enabledAt}},
				/* Life Ended            */ {Offer: Offer{Type: jentity.OfferTypeLife}, Contract: &Contract{EndDate: new(jdate.NewFromTime(time.Now().Add(-48 * time.Hour)))}},
				/* Open Ended            */ {Offer: Offer{Type: jentity.OfferTypeOpen, EnabledAt: &enabledAt, DisabledAt: &disabledAt}},
			},
		},
		{
			name: "full potential - churn today",
			original: []OfferWithContract{
				/* Open live                     */ {Offer: Offer{Type: jentity.OfferTypeOpen, EnabledAt: &enabledAt}},
				/* Life live futur churn (today) */ {Offer: Offer{Type: jentity.OfferTypeLife}, Contract: &Contract{EndDate: new(jdate.Today())}},
				/* Life Ended                    */ {Offer: Offer{Type: jentity.OfferTypeLife}, Contract: &Contract{EndDate: new(jdate.NewFromTime(time.Now().Add(-48 * time.Hour)))}},
				/* life onboarding               */ {Offer: Offer{Type: jentity.OfferTypeLife}},
				/* life live                     */ {Offer: Offer{Type: jentity.OfferTypeLife}, Contract: &Contract{}},
			},
			result: []OfferWithContract{
				/* Life live futur churn (today) */ {Offer: Offer{Type: jentity.OfferTypeLife}, Contract: &Contract{EndDate: new(jdate.Today())}},
				/* life live             */ {Offer: Offer{Type: jentity.OfferTypeLife}, Contract: &Contract{}},
				/* life onboarding       */ {Offer: Offer{Type: jentity.OfferTypeLife}},
				/* Open live             */ {Offer: Offer{Type: jentity.OfferTypeOpen, EnabledAt: &enabledAt}},
				/* Life Ended            */ {Offer: Offer{Type: jentity.OfferTypeLife}, Contract: &Contract{EndDate: new(jdate.NewFromTime(time.Now().Add(-48 * time.Hour)))}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OfferWithContractSortBy(FrontendDisplay).Sort(tt.original)

			assert.Equal(t, tt.result, tt.original)
		})
	}
}

func Test_Contract_IsActive(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		c    Contract
		d    jdate.Date
		want bool
	}{
		"active, first day of the contract, end date nil": {
			c: Contract{
				StartDate: jdate.New(2024, time.April, 2),
			},
			d:    jdate.New(2024, time.April, 2),
			want: true,
		},
		"active, end date known": {
			c: Contract{
				StartDate: jdate.New(2024, time.April, 2),
				EndDate:   new(jdate.New(2024, time.April, 3)),
			},
			d:    jdate.New(2024, time.April, 2),
			want: true,
		},
		"not active, before contract start": {
			c: Contract{
				StartDate: jdate.New(2024, time.April, 2),
			},
			d:    jdate.New(2024, time.April, 1),
			want: false,
		},
		"not active, after contract end": {
			c: Contract{
				StartDate: jdate.New(2024, time.April, 2),
				EndDate:   new(jdate.New(2024, time.April, 3)),
			},
			d:    jdate.New(2024, time.April, 4),
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.IsActive(tt.d))
		})
	}
}

func Test_ListUsersInfoParams_AddToQuery(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		params ListUsersInfoParams
		uri    url.URL
		want   url.URL
	}{
		"all fields set": {
			params: ListUsersInfoParams{
				QueryPagination: jdb.QueryPagination{
					Offset: 1,
					Limit:  10,
				},
				Search:          "test",
				Offer:           ListUsersInfoParamsOffer("blue"),
				Status:          ListUsersInfoParamsStatus("live"),
				Region:          jentity.RegionUK,
				SalesOwnerEmail: "sales@example.com",
				CSMOwnerEmail:   "csm@example.com",
			},
			uri: url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/users",
			},
			want: url.URL{
				Scheme:   "https",
				Host:     "example.com",
				Path:     "/users",
				RawQuery: "csm_owner_email=csm%40example.com&limit=10&offer=blue&offset=1&region=uk&sales_owner_email=sales%40example.com&search=test&status=live",
			},
		},
		"no optional fields set": {
			params: ListUsersInfoParams{
				QueryPagination: jdb.QueryPagination{
					Offset: 1,
					Limit:  10,
				},
			},
			uri: url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/users",
			},
			want: url.URL{
				Scheme:   "https",
				Host:     "example.com",
				Path:     "/users",
				RawQuery: "limit=10&offset=1",
			},
		},
		"some optional fields set": {
			params: ListUsersInfoParams{
				QueryPagination: jdb.QueryPagination{
					Offset: 1,
					Limit:  10,
				},
				Search: "test",
				Offer:  ListUsersInfoParamsOffer("blue"),
			},
			uri: url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/users",
			},
			want: url.URL{
				Scheme:   "https",
				Host:     "example.com",
				Path:     "/users",
				RawQuery: "limit=10&offer=blue&offset=1&search=test",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.params.AddToQuery(tt.uri)
			assert.Equal(t, tt.want.String(), got.String())
		})
	}
}
