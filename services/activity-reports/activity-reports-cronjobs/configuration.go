package main

import (
	"os"

	"github.com/Freelance-launchpad/backend-interview/common/jdb"
	"github.com/Freelance-launchpad/backend-interview/common/jutils"
)

type configuration struct {
	Database jdb.Config `yaml:"database"`

	AuthClientSecret string
	AuthURL          jutils.URL `yaml:"auth_url"`
	UsersURL         jutils.URL `yaml:"users_url"`
	GotenbergURL     jutils.URL `yaml:"gotenberg_url"`
}

func (c *configuration) Secrets() error {
	c.Database.Password = os.Getenv("DATABASE_PASSWORD")

	c.AuthClientSecret = os.Getenv("AUTH_CLIENT_SECRET")

	return nil
}
