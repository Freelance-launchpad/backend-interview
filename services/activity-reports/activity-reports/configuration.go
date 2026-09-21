package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/Freelance-launchpad/backend-interview/common/jdb"
	"github.com/Freelance-launchpad/backend-interview/common/jgin"
	"github.com/Freelance-launchpad/backend-interview/common/jutils"
)

type configuration struct {
	Database         jdb.Config  `yaml:"database"`
	Authorization    jgin.Config `yaml:"authorization"`
	JWTPublicKey     ed25519.PublicKey
	BannedJTIs       []string `yaml:"banned_jtis"`
	AuthClientSecret string
	AuthURL          jutils.URL `yaml:"auth_url"`
	UsersURL         jutils.URL `yaml:"users_url"`
	ExpensesURL      jutils.URL `yaml:"expenses_url"`
	PayrollURL       jutils.URL `yaml:"payroll_url"`
}

func (c *configuration) Secrets() error {
	c.Database.Password = os.Getenv("DATABASE_PASSWORD")

	publicKeyB64 := os.Getenv("JWT_PUBLIC_KEY")
	publicBytes, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return fmt.Errorf("failed to decode JWT public key: %w", err)
	}
	c.JWTPublicKey = ed25519.PublicKey(publicBytes)

	c.AuthClientSecret = os.Getenv("AUTH_CLIENT_SECRET")

	return nil
}
