package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jkafka/internal"
	"github.com/Freelance-launchpad/backend-interview/common/jutils"
)

type Config struct {
	URL jutils.URL `yaml:"url"`
}

type client struct {
	jhttp.Doer
	url url.URL
}

func New(httpClient jhttp.Doer, conf Config) internal.RegistryClient {
	return &client{
		Doer: httpClient,
		url:  *conf.URL.URL,
	}
}

// registerRequest is the payload sent to the Schema Registry.
type registerRequest struct {
	Schema     string `json:"schema"`
	SchemaType string `json:"schemaType"`
}

// registerResponse is the response from the Schema Registry.
type registerResponse struct {
	ID int32 `json:"id"`
}

// RegisterSchema registers a JSON Schema under the given subject and returns the schema ID.
// If the schema is already registered, the existing ID is returned (idempotent).
func (c *client) RegisterSchema(ctx context.Context, subject, schema string) (_ int32, err error) {
	defer jerror.Wrap(&err, "with subject", subject)

	uri := c.url
	uri.Path = fmt.Sprintf("/subjects/%s/versions", subject)

	body, err := json.Marshal(registerRequest{
		Schema:     schema,
		SchemaType: "JSON",
	})
	if err != nil {
		return 0, fmt.Errorf("marshal register request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), bytes.NewReader(body))
	if err != nil {
		return 0, fmt.Errorf("build register request: %w", err)
	}
	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")

	resp, err := c.Do(req)
	if err != nil {
		return 0, fmt.Errorf("register schema: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			responseBody = fmt.Appendf(nil, "error reading response body: %v", err)
		}
		return 0, jhttp.APIError{
			Err:        fmt.Errorf("bad status, expected: %d - with body %s", http.StatusOK, string(responseBody)),
			StatusCode: resp.StatusCode,
		}
	}

	var result registerResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decode register response: %w", err)
	}

	return result.ID, nil
}

// GetSchema fetches the raw schema string for a given schema ID.
func (c *client) GetSchema(ctx context.Context, id int32) (_ string, err error) {
	defer jerror.Wrap(&err, "with id", strconv.FormatInt(int64(id), 10))

	uri := c.url
	uri.Path = fmt.Sprintf("/schemas/ids/%d", id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return "", fmt.Errorf("build get schema request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("get schema %d: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			responseBody = fmt.Appendf(nil, "error reading response body: %v", err)
		}
		return "", jhttp.APIError{
			Err:        fmt.Errorf("bad status, expected: %d - with body %s", http.StatusOK, string(responseBody)),
			StatusCode: resp.StatusCode,
		}
	}

	var result struct {
		Schema string `json:"schema"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode schema response: %w", err)
	}

	return result.Schema, nil
}
