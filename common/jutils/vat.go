package jutils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/teamwork/vat/v3"
)

type VATNumber string

func (s *VATNumber) UnmarshalJSON(data []byte) error {
	var vatString string
	if err := json.Unmarshal(data, &vatString); err != nil {
		return err
	}

	vatString = strings.ReplaceAll(vatString, " ", "")

	if err := vat.ValidateFormat(vatString); err != nil {
		return fmt.Errorf("invalid VAT number: %w", err)
	}

	*s = VATNumber(vatString)
	return nil
}
