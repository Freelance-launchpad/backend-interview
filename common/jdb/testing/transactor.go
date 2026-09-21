package testing

import (
	"context"
	"database/sql"

	"github.com/Thiht/transactor/stdlib"
)

func NewFakeTransactor() fakeTransactor {
	return fakeTransactor{}
}

func NewFakeDBGetter(db *sql.DB) stdlib.DBGetter {
	return func(context.Context) stdlib.DB {
		return db
	}
}

type fakeTransactorKey struct{}

type fakeTransactor struct{}

func (f fakeTransactor) WithinTransaction(ctx context.Context, txFunc func(context.Context) error) error {
	return txFunc(Context(ctx))
}

// Context adds the fake transactor key to the provided context.
// This can be used to simulate a WithinTransaction context, which is useful when comparing contexts in tests.
func Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, fakeTransactorKey{}, nil)
}
