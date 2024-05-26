package pg

import (
	"backend/internal/repo"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var _ repo.AuthRepo = (*AuthRepo)(nil)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

const authRepoSetRefreshTokenQuery = `
insert into public.user_token(user_id, refresh)
values ($1, $2)
on conflict (user_id) do update set refresh = $2
`

func (r *AuthRepo) SetRefreshToken(
	ctx context.Context,
	userId int64,
	token string,
) error {
	if _, err := r.db.ExecContext(ctx, authRepoSetRefreshTokenQuery, userId, token); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}
	return nil
}

const authRepoVerifyRefreshTokenQuery = `
select 1
from public.user_token
where user_id = $1 and refresh = $2
`

func (r *AuthRepo) VerifyRefreshToken(
	ctx context.Context,
	userId int64,
	token string,
) (bool, error) {
	var ok int8
	if err := r.db.GetContext(ctx, &ok, authRepoVerifyRefreshTokenQuery, userId, token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, repo.ErrNotFound
		}
		return false, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return true, nil
}
