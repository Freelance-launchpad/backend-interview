// SQL table creation example:
//
//	CREATE TABLE oauth_tokens (
//		token_name TEXT PRIMARY KEY NOT NULL
//		, access_token TEXT NOT NULL
//		, refresh_token TEXT NOT NULL
//		, expires_at TIMESTAMPTZ NOT NULL
//		, updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
//	);
package oauth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Freelance-launchpad/backend-interview/common/jdb"
	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
)

type Store interface {
	BeginLocked(ctx context.Context) (jdb.RunnerTx, error)
	GetToken(ctx context.Context, tx jdb.Runner, tokenName string) (Token, error)
	UpdateToken(ctx context.Context, tx jdb.Runner, token Token) error
}

type store struct {
	*sql.DB
}

func NewStore(db *sql.DB) Store {
	return &store{
		DB: db,
	}
}

func (s *store) BeginLocked(ctx context.Context) (_ jdb.RunnerTx, err error) {
	defer jerror.Wrap(&err)

	tx, err := s.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, "LOCK TABLE oauth_tokens")
	if err != nil {
		return nil, err
	}

	return tx, nil
}

const queryGetToken = `
	SELECT
		token_name,
		access_token,
		refresh_token,
		expires_at
	FROM oauth_tokens
	WHERE token_name = $1
	LIMIT 1;
`

func (s *store) GetToken(ctx context.Context, tx jdb.Runner, tokenName string) (_ Token, err error) {
	defer jerror.Wrap(&err, "with token_name", tokenName)

	if len(tokenName) == 0 {
		return Token{}, fmt.Errorf("invalid token name")
	}

	runner := tx
	if tx == nil {
		runner = s.DB
	}

	var token Token
	err = runner.QueryRowContext(ctx, queryGetToken, tokenName).Scan(
		&token.Name,
		&token.AccessToken,
		&token.RefreshToken,
		&token.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Token{}, jerror.NewNotFoundError("TOKEN_NOT_FOUND")
		}
		return Token{}, err
	}
	return token, nil
}

const queryUpdateToken = `
	UPDATE oauth_tokens
	SET access_token = $2
		, refresh_token = $3
		, expires_at = $4
		, updated_at = NOW()
	WHERE token_name = $1;
`

func (s *store) UpdateToken(ctx context.Context, tx jdb.Runner, token Token) (err error) {
	defer jerror.Wrap(&err, "with token_name", token.Name)

	runner := tx
	if tx == nil {
		runner = s.DB
	}

	res, err := runner.ExecContext(
		ctx,
		queryUpdateToken,
		token.Name,
		token.AccessToken,
		token.RefreshToken,
		token.ExpiresAt,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return jerror.NewNotFoundError("TOKEN_NOT_FOUND")
	}

	return nil
}
