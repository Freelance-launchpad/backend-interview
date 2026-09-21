package jdb

import (
	"errors"

	"github.com/lib/pq"
)

const (
	PostgresUniqueViolationCode     = "23505"
	PostgresInvalidForeignKey       = "42830"
	PostgresForeignKeyViolationCode = "23503"
)

var (
	ErrUniqueViolation = &pq.Error{
		Code: PostgresUniqueViolationCode,
	}
	ErrInvalidForeignKey = &pq.Error{
		Code: PostgresInvalidForeignKey,
	}
)

func IsErrorUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == PostgresUniqueViolationCode
	}

	return false
}

func IsErrorInvalidForeignKey(err error) bool {
	if err == nil {
		return false
	}

	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == PostgresInvalidForeignKey
	}

	return false
}

func IsErrorForeignKeyViolationOnConstraint(err error, constraint string) bool {
	if pqErr, ok := errors.AsType[*pq.Error](err); ok {
		return pqErr.Code == PostgresForeignKeyViolationCode && pqErr.Constraint == constraint
	}

	return false
}
