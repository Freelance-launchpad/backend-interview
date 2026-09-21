package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/v2"
	"github.com/Freelance-launchpad/backend-interview/services/expenses/pkg"
	"github.com/google/uuid"
)

const (
	pathV3 = "/expenses-api/v3"
	pathV4 = "/expenses-api/v4"
)

type Client interface {
	GetSalaryExpenses(ctx context.Context, offerID uuid.UUID) ([]pkg.SalaryExpense, error)
	ListExpenses(ctx context.Context, filter pkg.ExpenseFilter, withAttachments bool) (pkg.ListExpensesResponse, error)
	LockExpenses(ctx context.Context, offerID uuid.UUID, expensesIDs uuid.UUIDs) error
	UnlockExpenses(ctx context.Context, offerID uuid.UUID) error
	CompleteAndUnlockExpenses(ctx context.Context, offerID uuid.UUID, expenses []pkg.RefundedExpense) error
	SetExpensesNonRefundable(ctx context.Context, params []pkg.NonRefundableParams) error
	ListCategories(ctx context.Context) ([]pkg.Category, error)
}

type client struct {
	jhttp.Doer
	url url.URL
}

func New(httpClient jhttp.Doer, uri url.URL) Client {
	return &client{
		Doer: httpClient,
		url:  uri,
	}
}

func (c *client) GetSalaryExpenses(ctx context.Context, offerID uuid.UUID) ([]pkg.SalaryExpense, error) {
	const errMsg = "client.GetSalaryExpenses has failed"

	uri := c.url
	uri.Path = path.Join(pathV3, "admin/expenses/offers", offerID.String(), "salary")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: http.NewRequest err: %w", errMsg, err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: client.Do err: %w", errMsg, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.ErrorContext(ctx, errMsg+": unable to close response", slog.Any("error", err))
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		responseBody = fmt.Appendf(nil, "error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, jhttp.APIError{
			Err:        fmt.Errorf("%s: bad status, expected: %d, received: %d - with body %s", errMsg, http.StatusOK, resp.StatusCode, string(responseBody)),
			StatusCode: resp.StatusCode,
		}
	}

	var expenses []pkg.SalaryExpense
	if err := json.Unmarshal(responseBody, &expenses); err != nil {
		return nil, fmt.Errorf("%s: json.Unmarshal err: %w", errMsg, err)
	}

	return expenses, nil
}

func (c *client) ListExpenses(ctx context.Context, filter pkg.ExpenseFilter, withAttachments bool) (pkg.ListExpensesResponse, error) {
	const errMsg = "client.ListExpenses has failed"

	uri := c.url
	uri.Path = path.Join(pathV3, "admin/expenses")
	values := uri.Query()
	err := jhttp.AddStructInValues(values, filter)
	if err != nil {
		return pkg.ListExpensesResponse{}, fmt.Errorf("%s: AddStructInValues err: %w", errMsg, err)
	}
	values.Add("enriched_attachments", strconv.FormatBool(withAttachments))
	values.Add("with_auto_validation_result", "false")

	uri.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return pkg.ListExpensesResponse{}, fmt.Errorf("%s: http.NewRequest err: %w", errMsg, err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return pkg.ListExpensesResponse{}, fmt.Errorf("%s: client.Do err: %w", errMsg, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.ErrorContext(ctx, errMsg+": unable to close response", slog.Any("error", err))
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		responseBody = fmt.Appendf(nil, "error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return pkg.ListExpensesResponse{}, jhttp.APIError{
			Err:        fmt.Errorf("%s: bad status, expected: %d, received: %d - with body %s", errMsg, http.StatusOK, resp.StatusCode, string(responseBody)),
			StatusCode: resp.StatusCode,
		}
	}

	var expenses pkg.ListExpensesResponse
	if err := json.Unmarshal(responseBody, &expenses); err != nil {
		return pkg.ListExpensesResponse{}, fmt.Errorf("%s: json.Unmarshal err: %w", errMsg, err)
	}

	return expenses, nil
}

func (c *client) LockExpenses(ctx context.Context, offerID uuid.UUID, expensesIDs uuid.UUIDs) error {
	const errMsg = "client.LockExpenses has failed"

	uri := c.url
	uri.Path = path.Join(pathV3, "admin/expenses/lock")

	query := pkg.LockExpensesQuery{
		OfferID:     offerID,
		ExpensesIDs: expensesIDs,
	}
	body, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("%s: json.Marshal err: %w", errMsg, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("%s: http.NewRequest err: %w", errMsg, err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("%s: client.Do err: %w", errMsg, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.ErrorContext(ctx, errMsg+": unable to close response", slog.Any("error", err))
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		responseBody = fmt.Appendf(nil, "error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return jhttp.APIError{
			Err:        fmt.Errorf("%s: bad status, expected: %d, received: %d - with body %s", errMsg, http.StatusNoContent, resp.StatusCode, string(responseBody)),
			StatusCode: resp.StatusCode,
		}
	}

	return nil
}

func (c *client) UnlockExpenses(ctx context.Context, offerID uuid.UUID) (err error) {
	defer jerror.Wrap(&err, "with offerID", offerID.String())

	uri := c.url
	uri.Path = path.Join(pathV3, "admin/expenses/offers/"+offerID.String()+"/unlock")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest err: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do err: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		responseBody = fmt.Appendf(nil, "error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return jhttp.APIError{
			Err:        fmt.Errorf("bad status, expected: %d, received: %d - with body %s", http.StatusNoContent, resp.StatusCode, string(responseBody)),
			StatusCode: resp.StatusCode,
		}
	}

	return nil
}

func (c *client) CompleteAndUnlockExpenses(ctx context.Context, offerID uuid.UUID, expenses []pkg.RefundedExpense) (err error) {
	defer jerror.Wrap(&err)

	uri := c.url
	uri.Path = path.Join(pathV4, fmt.Sprintf("admin/%s/expenses/complete", offerID))

	body, err := json.Marshal(pkg.RefundedExpenses{
		Expenses: expenses,
	})
	if err != nil {
		return fmt.Errorf("json.Marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http.NewRequest err: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do err: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.ErrorContext(ctx, "unable to close response", slog.Any("error", err))
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		responseBody = fmt.Appendf(nil, "error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return jhttp.APIError{
			Err:        fmt.Errorf("bad status, expected: %d, received: %d - with body %s", http.StatusNoContent, resp.StatusCode, string(responseBody)),
			StatusCode: resp.StatusCode,
		}
	}

	return nil
}

func (c *client) SetExpensesNonRefundable(ctx context.Context, params []pkg.NonRefundableParams) error {
	const errMsg = "client.SetExpensesNonRefundable has failed"

	uri := c.url
	uri.Path = path.Join(pathV3, "admin/expenses/non_refundable")

	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("%s: json.Marshal error: %w", errMsg, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("%s: http.NewRequest err: %w", errMsg, err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("%s: client.Do err: %w", errMsg, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.ErrorContext(ctx, errMsg+": unable to close response", slog.Any("error", err))
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		responseBody = fmt.Appendf(nil, "error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return jhttp.APIError{
			Err:        fmt.Errorf("%s: bad status, expected: %d, received: %d - with body %s", errMsg, http.StatusNoContent, resp.StatusCode, string(responseBody)),
			StatusCode: resp.StatusCode,
		}
	}

	return nil
}

func (c *client) ListCategories(ctx context.Context) (categories []pkg.Category, err error) {
	defer jerror.Wrap(&err)

	uri := c.url
	uri.Path = path.Join(pathV3, "admin/categories")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest err: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do err: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		responseBody = fmt.Appendf(nil, "error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, jhttp.APIError{
			Err:        fmt.Errorf("bad status, expected: %d, received: %d - with body %s", http.StatusOK, resp.StatusCode, string(responseBody)),
			StatusCode: resp.StatusCode,
		}
	}

	var result []pkg.Category
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("json.Unmarshal err: %w", err)
	}

	return result, nil
}
