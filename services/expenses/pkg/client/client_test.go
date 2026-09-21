package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jhttp/httpmock"
	"github.com/Freelance-launchpad/backend-interview/services/expenses/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Client_GetSalaryExpenses(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	offerID := uuid.New()

	path := "/expenses-api/v3/admin/expenses/offers/" + offerID.String() + "/salary"

	expenses := []pkg.SalaryExpense{
		{
			ID:       uuid.New(),
			Category: pkg.ExpenseCategoryNTICGear,
		},
	}

	tests := []struct {
		name   string
		client httpmock.Client
		output []pkg.SalaryExpense
		err    error
	}{
		{
			name: "nominal",
			client: *httpmock.New(t).
				WithCall(
					http.MethodGet,
					path,
					httpmock.ReturnStatus(http.StatusOK),
					httpmock.ReturnBody(func() string {
						body, err := json.Marshal(expenses)
						require.NoError(t, err)
						return string(body)
					}()),
				),
			output: expenses,
		},
		{
			name: "Do returns an error",
			client: *httpmock.New(t).
				WithCall(
					http.MethodGet,
					path,
					httpmock.ReturnError(assert.AnError),
				),
			err: assert.AnError,
		},
		{
			name: "expenses returns an error",
			client: *httpmock.New(t).
				WithCall(
					http.MethodGet,
					path,
					httpmock.ReturnStatus(http.StatusInternalServerError),
					httpmock.ReturnBody("an error"),
				),
			err: errors.New("[500]: client.GetSalaryExpenses has failed: bad status, expected: 200, received: 500 - with body an error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := client{Doer: &tt.client}

			output, err := client.GetSalaryExpenses(ctx, offerID)

			assert.Equal(t, tt.output, output)
			if tt.err != nil {
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Client_ListExpenses(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	path := "/expenses-api/v3/admin/expenses"

	offerID1, offerID2 := uuid.New(), uuid.New()

	expenses := pkg.ListExpensesResponse{
		Expenses: []pkg.Expense{
			{
				ID:       uuid.New(),
				Category: pkg.ExpenseCategoryNTICGear,
			},
		},
		Count: 2,
	}

	tests := []struct {
		name            string
		filter          pkg.ExpenseFilter
		withAttachments bool
		client          httpmock.Client
		output          pkg.ListExpensesResponse
		err             error
	}{
		{
			name: "nominal - no filter",
			client: *httpmock.New(t).
				WithCall(
					http.MethodGet,
					path,
					httpmock.ExpectQueryParam("enriched_attachments", "false"),
					httpmock.ExpectQueryParam("with_auto_validation_result", "false"),
					httpmock.ReturnStatus(http.StatusOK),
					httpmock.ReturnBody(func() string {
						body, err := json.Marshal(expenses)
						require.NoError(t, err)
						return string(body)
					}()),
				),
			output: expenses,
		},
		{
			name: "nominal - with filters",
			filter: pkg.ExpenseFilter{
				OfferIDs:       []uuid.UUID{offerID1, offerID2},
				MaxExpenseDate: new(time.Date(2024, time.August, 16, 15, 47, 00, 0, time.UTC)),
			},
			withAttachments: true,
			client: *httpmock.New(t).
				WithCall(
					http.MethodGet,
					path,
					httpmock.ExpectQueryParam("enriched_attachments", "true"),
					httpmock.ExpectQueryParam("with_auto_validation_result", "false"),
					httpmock.ExpectQueryParam("offer_ids", offerID1.String()),
					httpmock.ExpectQueryParam("offer_ids", offerID2.String()),
					httpmock.ExpectQueryParam("expense_date", "2024-08-16T15:47:00Z"),
					httpmock.ReturnStatus(http.StatusOK),
					httpmock.ReturnBody(func() string {
						body, err := json.Marshal(expenses)
						require.NoError(t, err)
						return string(body)
					}()),
				),
			output: expenses,
		},
		{
			name: "Do returns an error",
			client: *httpmock.New(t).
				WithCall(
					http.MethodGet,
					path,
					httpmock.ReturnError(assert.AnError),
				),
			err: assert.AnError,
		},
		{
			name: "expenses returns an error",
			client: *httpmock.New(t).
				WithCall(
					http.MethodGet,
					path,
					httpmock.ReturnStatus(http.StatusInternalServerError),
					httpmock.ReturnBody("an error"),
				),
			err: errors.New("[500]: client.ListExpenses has failed: bad status, expected: 200, received: 500 - with body an error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := client{Doer: &tt.client}

			output, err := client.ListExpenses(ctx, tt.filter, tt.withAttachments)

			assert.Equal(t, tt.output, output)
			if tt.err != nil {
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Client_LockExpenses(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	path := "/expenses-api/v3/admin/expenses/lock"

	offerID := uuid.New()
	expensesIDs := uuid.UUIDs{uuid.New(), uuid.New()}

	body, err := json.Marshal(pkg.LockExpensesQuery{
		OfferID:     offerID,
		ExpensesIDs: expensesIDs,
	})
	require.NoError(t, err)
	expectedBody := string(body)

	tests := []struct {
		name   string
		client httpmock.Client
		err    error
	}{
		{
			name: "nominal",
			client: *httpmock.New(t).
				WithCall(
					http.MethodPost,
					path,
					httpmock.ExpectJSON(expectedBody),
					httpmock.ReturnStatus(http.StatusNoContent),
				),
		},
		{
			name: "Do returns an error",
			client: *httpmock.New(t).
				WithCall(
					http.MethodPost,
					path,
					httpmock.ExpectJSON(expectedBody),
					httpmock.ReturnError(assert.AnError),
				),
			err: assert.AnError,
		},
		{
			name: "expenses returns an error",
			client: *httpmock.New(t).
				WithCall(
					http.MethodPost,
					path,
					httpmock.ExpectJSON(expectedBody),
					httpmock.ReturnStatus(http.StatusInternalServerError),
					httpmock.ReturnBody("an error"),
				),
			err: errors.New("[500]: client.LockExpenses has failed: bad status, expected: 204, received: 500 - with body an error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := client{Doer: &tt.client}

			err := client.LockExpenses(ctx, offerID, expensesIDs)

			if tt.err != nil {
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Client_UnlockExpenses(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	offerID := uuid.New()
	path := "/expenses-api/v3/admin/expenses/offers/" + offerID.String() + "/unlock"

	tests := []struct {
		name   string
		client httpmock.Client
		err    error
	}{
		{
			name: "nominal",
			client: *httpmock.New(t).
				WithCall(
					http.MethodPost,
					path,
					httpmock.ReturnStatus(http.StatusNoContent),
				),
		},
		{
			name: "Do returns an error",
			client: *httpmock.New(t).
				WithCall(
					http.MethodPost,
					path,
					httpmock.ReturnError(assert.AnError),
				),
			err: assert.AnError,
		},
		{
			name: "expenses returns an error",
			client: *httpmock.New(t).
				WithCall(
					http.MethodPost,
					path,
					httpmock.ReturnStatus(http.StatusInternalServerError),
					httpmock.ReturnBody("an error"),
				),
			err: errors.New("[500]: bad status, expected: 204, received: 500 - with body an error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := client{Doer: &tt.client}

			err := client.UnlockExpenses(ctx, offerID)

			if tt.err != nil {
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Client_SetExpensesNonRefundable(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	path := "/expenses-api/v3/admin/expenses/non_refundable"

	params := []pkg.NonRefundableParams{
		{OfferID: uuid.New()},
	}

	body, err := json.Marshal(params)
	require.NoError(t, err)
	expectedBody := string(body)

	tests := []struct {
		name   string
		client httpmock.Client
		err    error
	}{
		{
			name: "nominal",
			client: *httpmock.New(t).
				WithCall(
					http.MethodPut,
					path,
					httpmock.ExpectJSON(expectedBody),
					httpmock.ReturnStatus(http.StatusNoContent),
				),
		},
		{
			name: "Do returns an error",
			client: *httpmock.New(t).
				WithCall(
					http.MethodPut,
					path,
					httpmock.ExpectJSON(expectedBody),
					httpmock.ReturnError(assert.AnError),
				),
			err: assert.AnError,
		},
		{
			name: "expenses returns an error",
			client: *httpmock.New(t).
				WithCall(
					http.MethodPut,
					path,
					httpmock.ExpectJSON(expectedBody),
					httpmock.ReturnStatus(http.StatusInternalServerError),
					httpmock.ReturnBody("an error"),
				),
			err: errors.New("[500]: client.SetExpensesNonRefundable has failed: bad status, expected: 204, received: 500 - with body an error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := client{Doer: &tt.client}

			err := client.SetExpensesNonRefundable(ctx, params)

			if tt.err != nil {
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Client_ListCategories(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	path := "/expenses-api/v3/admin/categories"

	categories := []pkg.Category{
		{
			ID:                pkg.ExpenseCategoryPersonalMeal,
			Icon:              "meal",
			Name:              "Repas personnel",
			Selectable:        true,
			ReimbursementRate: 100,
		},
		{
			ID:                pkg.ExpenseCategoryNTICGear,
			Icon:              "computer",
			Name:              "Matériel NTIC",
			Selectable:        true,
			ReimbursementRate: 100,
		},
	}

	tests := []struct {
		name   string
		client httpmock.Client
		output []pkg.Category
		err    error
	}{
		{
			name: "nominal",
			client: *httpmock.New(t).
				WithCall(
					http.MethodGet,
					path,
					httpmock.ReturnStatus(http.StatusOK),
					httpmock.ReturnBody(func() string {
						body, err := json.Marshal(categories)
						require.NoError(t, err)
						return string(body)
					}()),
				),
			output: categories,
		},
		{
			name: "Do returns an error",
			client: *httpmock.New(t).
				WithCall(
					http.MethodGet,
					path,
					httpmock.ReturnError(assert.AnError),
				),
			err: assert.AnError,
		},
		{
			name: "API returns an error",
			client: *httpmock.New(t).
				WithCall(
					http.MethodGet,
					path,
					httpmock.ReturnStatus(http.StatusInternalServerError),
					httpmock.ReturnBody("an error"),
				),
			err: errors.New("[500]: bad status, expected: 200, received: 500 - with body an error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := client{Doer: &tt.client}

			output, err := client.ListCategories(ctx)

			assert.Equal(t, tt.output, output)
			if tt.err != nil {
				assert.ErrorContains(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
