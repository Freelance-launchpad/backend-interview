package jgin

import (
	"database/sql/driver"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

// UUID is a wrapper around uuid.UUID to implement the gin binding interface.
//
// FIXME: This should be modified as `type UUID uuid.UUID` when https://github.com/gin-gonic/gin/issues/4265 is resolved.
type UUID struct {
	uuid.UUID
}

var _ binding.BindUnmarshaler = (*UUID)(nil)

func (u *UUID) UnmarshalParam(param string) error {
	id, err := uuid.Parse(param)
	if err != nil {
		return err
	}

	*u = UUID{id}
	return nil
}

var _ driver.Valuer = (*UUID)(nil)

func (u UUID) Value() (driver.Value, error) {
	return u.UUID.Value()
}

func (u *UUID) ToUUIDPtr() *uuid.UUID {
	if u == nil {
		return nil
	}
	return &u.UUID
}

var _ binding.BindUnmarshaler = (*jdate.Date)(nil)
