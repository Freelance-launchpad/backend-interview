package jutils

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"

	"github.com/almerlucke/go-iban/iban"
	"gopkg.in/yaml.v3"
)

type IBAN iban.IBAN

func ParseIBAN(ibanString string) (IBAN, error) {
	parsed, err := iban.NewIBAN(ibanString)
	if err != nil {
		return IBAN{}, err
	}

	return IBAN(*parsed), nil
}

func (i IBAN) String() string {
	return i.Code
}

func (i IBAN) PrettyString() string {
	var prettyIBAN bytes.Buffer
	for i, r := range i.Code {
		if i > 0 && i%4 == 0 {
			prettyIBAN.WriteRune(' ')
		}
		prettyIBAN.WriteRune(r)
	}
	return prettyIBAN.String()
}

var _ json.Marshaler = IBAN{}

func (i IBAN) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Code)
}

var _ json.Unmarshaler = &IBAN{}

func (i *IBAN) UnmarshalJSON(data []byte) error {
	var ibanString string
	if err := json.Unmarshal(data, &ibanString); err != nil {
		return err
	}

	parsed, err := iban.NewIBAN(strings.TrimSpace(ibanString))
	if err != nil {
		return err
	}

	*i = IBAN(*parsed)
	return nil
}

var _ yaml.Unmarshaler = &IBAN{}

func (i *IBAN) UnmarshalYAML(node *yaml.Node) error {
	var ibanString string
	if err := node.Decode(&ibanString); err != nil {
		return err
	}

	parsed, err := iban.NewIBAN(strings.TrimSpace(ibanString))
	if err != nil {
		return err
	}

	*i = IBAN(*parsed)
	return err
}

var _ driver.Valuer = IBAN{}

func (i IBAN) Value() (driver.Value, error) {
	return i.Code, nil
}

var _ sql.Scanner = &IBAN{}

func (i *IBAN) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		if v == "" {
			return nil
		}
		parsed, err := iban.NewIBAN(v)
		if err != nil {
			return err
		}

		*i = IBAN(*parsed)
		return nil

	case []byte:
		if len(v) == 0 {
			return nil
		}
		parsed, err := iban.NewIBAN(string(v))
		if err != nil {
			return err
		}

		*i = IBAN(*parsed)
		return nil

	default:
		return errors.New("invalid type for IBAN")
	}
}

// BBANFR represents the BBAN part of a FR IBAN.
type BBANFR struct {
	Raw                string
	BankCode           string
	BranchCode         string
	BankAccountNumber  string
	NationalCheckDigit string
}

func NewBBANFR(iban IBAN) (BBANFR, error) {
	if iban.CountryCode != "FR" {
		return BBANFR{}, errors.New("IBAN is not from France")
	}

	bban := iban.BBAN
	if len(bban) != 23 {
		return BBANFR{}, errors.New("invalid BBAN length for FR IBAN")
	}

	return BBANFR{
		Raw:                bban,
		BankCode:           bban[0:5],
		BranchCode:         bban[5:10],
		BankAccountNumber:  bban[10:21],
		NationalCheckDigit: bban[21:23],
	}, nil
}

type BBANUK struct {
	Raw           string
	BankCode      string
	SortCode      string
	AccountNumber string
}

func NewBBANUK(iban IBAN) (BBANUK, error) {
	if iban.CountryCode != "GB" {
		return BBANUK{}, errors.New("IBAN is not from the UK")
	}

	bban := iban.BBAN
	if len(bban) != 18 {
		return BBANUK{}, errors.New("invalid BBAN length for GB IBAN")
	}

	return BBANUK{
		Raw:           bban,
		BankCode:      bban[0:4],
		SortCode:      bban[4:10],
		AccountNumber: bban[10:18],
	}, nil
}
