package jgin

type Config struct {
	// ISS is the server issuer identifier.
	ISS string `yaml:"iss"`

	// JWKURL is a URL for the JWK resource.
	JWKURL string `yaml:"jwk_url"`

	// CustomFieldPrefix is the prefix for custom field.
	CustomFieldPrefix string `yaml:"custom_field_prefix"`

	// SkipAuth is field used to skip authentication verification during tests.
	SkipAuth bool `yaml:"skip_auth"`
}
