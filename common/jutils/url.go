package jutils

import (
	"encoding/json"
	"net/url"

	"gopkg.in/yaml.v3"
)

// URL is a custom type which allows to marshal and unmarshal url value from yaml and json format.
type URL struct {
	*url.URL
}

// MarshalJSON cast URL type into json string (during json marshal)
func (j URL) MarshalJSON() ([]byte, error) {
	if j.URL == nil {
		return []byte("null"), nil
	}
	return json.Marshal(j.String())
}

// UnmarshalJSON cast json string value into URL type (during a json unmarshal)
func (j *URL) UnmarshalJSON(data []byte) error {
	var urlString string
	if err := json.Unmarshal(data, &urlString); err != nil {
		return err
	}
	u, err := url.Parse(urlString)
	if err != nil {
		return err
	}
	j.URL = u
	return nil
}

var _ yaml.Unmarshaler = &URL{}

// UnmarshalYAML cast yaml string value into URL type (during a yaml unmarshal)
func (j *URL) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}
	url, err := url.Parse(s)
	j.URL = url
	return err
}

var _ yaml.Marshaler = URL{}

// MarshalYAML cast URL type into yaml string (during yaml marshal)
func (j URL) MarshalYAML() (any, error) {
	return j.String(), nil
}
