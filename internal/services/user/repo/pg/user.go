package pg

import (
	"backend/internal/models"
	"backend/internal/services/user/repo"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

const getByIdQuery = `
select 
    u.id, 
    u.group_id, 
    u.email, 
    u.password, 
    u.first_name, 
    u.last_name, 
    u.middle_name, 
    u.created_at, 
    u.updated_at
from public.user u
where u.id = $1
`

func (r *UsersRepo) GetById(
	ctx context.Context,
	id int64,
) (models.User, error) {
	var u user
	if err := r.db.GetContext(ctx, &u, getByIdQuery, id); err != nil {
		return models.User{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return u.toServiceModel(), nil
}

const getByCredentialsQuery = `
select 
    u.id, 
    u.group_id, 
    u.email, 
    u.password, 
    u.first_name, 
    u.last_name, 
    u.middle_name, 
    u.created_at, 
    u.updated_at
from public.user u
where u.email = $1 and u.password = $2
`

func (r *UsersRepo) GetByCredentials(
	ctx context.Context,
	opts repo.UsersRepoGetByCredentialsOpts,
) (models.User, error) {
	var u user
	if err := r.db.GetContext(ctx, &u, getByCredentialsQuery, opts.Email, opts.Password); err != nil {
		return models.User{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return u.toServiceModel(), nil
}

const createQuery = `
insert into public.user (group_id, email, password, first_name, last_name, middle_name) 
values (:group_id, :email, :password, :first_name, :last_name, :middle_name)
`

func (r *UsersRepo) Create(
	ctx context.Context,
	opts repo.UsersRepoCreateOpts,
) (models.User, error) {
	_, err := r.db.NamedExecContext(ctx, createQuery, struct {
		GroupId    int64   `db:"group_id"`
		Email      string  `db:"email"`
		Password   string  `db:"password"`
		FirstName  string  `db:"first_name"`
		LastName   string  `db:"last_name"`
		MiddleName *string `db:"middle_name"`
	}{
		GroupId:    opts.GroupId,
		Email:      opts.Email,
		Password:   opts.Password,
		FirstName:  opts.FirstName,
		LastName:   opts.LastName,
		MiddleName: opts.MiddleName,
	})
	if err != nil {
		return models.User{}, fmt.Errorf("r.db.NamedExecContext: %w", err)
	}

	return models.User{}, nil
}
