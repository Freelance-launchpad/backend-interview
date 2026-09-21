package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/v2"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/pkg"
	"github.com/google/uuid"
)

type Client interface {
	GetRestDays(ctx context.Context, offerID uuid.UUID, period jdate.YearMonth) (pkg.RestDaysResponse, error)
}

const pathV3 = "/activity-reports/v3"

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

func (c *client) GetRestDays(
	ctx context.Context,
	offerID uuid.UUID,
	period jdate.YearMonth,
) (_ pkg.RestDaysResponse, err error) {

	defer jerror.Wrap(&err)

	uri := c.url
	uri.Path = path.Join(pathV3, "/admin/offers/", offerID.String(), "/work-schedules/rest-days")
	values := uri.Query()
	values.Set("year", strconv.Itoa(period.Year))
	values.Set("month", strconv.Itoa(int(period.Month)))
	uri.RawQuery = values.Encode()

	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return pkg.RestDaysResponse{}, fmt.Errorf("new request error: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return pkg.RestDaysResponse{}, fmt.Errorf("do error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			responseBody = fmt.Appendf(nil, "error reading response body: %s", err.Error())
		}
		return pkg.RestDaysResponse{}, jhttp.APIError{
			Err:        fmt.Errorf("bad status, expected: %d, with body: %s", http.StatusOK, responseBody),
			StatusCode: resp.StatusCode,
		}
	}

	var response pkg.RestDaysResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return pkg.RestDaysResponse{}, fmt.Errorf("decode error: %w", err)
	}

	return response, nil
}
