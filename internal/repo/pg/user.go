package pg

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"time"

	"github.com/jmoiron/sqlx"
)

type user struct {
	Id         int64      `db:"id"`
	GroupId    *int64     `db:"group_id"`
	RoleId     int64      `db:"role_id"`
	Email      string     `db:"email"`
	Password   string     `db:"password"`
	FirstName  string     `db:"first_name"`
	LastName   string     `db:"last_name"`
	MiddleName *string    `db:"middle_name"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

func (u user) toServiceModel() models.User {
	return models.User{
		Id:         u.Id,
		GroupId:    u.GroupId,
		RoleId:     u.RoleId,
		Email:      u.Email,
		Password:   u.Password,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		MiddleName: u.MiddleName,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

const usersRepoGetByIdQuery = `
select 
    u.id, 
    u.group_id, 
    u.role_id,
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
	if err := r.db.GetContext(ctx, &u, usersRepoGetByIdQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repo.ErrNotFound
		}
		return models.User{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return u.toServiceModel(), nil
}

const usersRepoGetListByRoleIdQuery = `
select *
from public.user
where role_id = $1
order by id desc 
limit $2
offset $3
`

func (r *UsersRepo) GetListByRoleId(
	ctx context.Context,
	opts repo.UsersRepoGetListByRoleIdOpts,
) ([]models.User, error) {
	var users []user
	if err := r.db.SelectContext(ctx, &users, usersRepoGetListByRoleIdQuery, opts.RoleId, opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("r.db.SelectContext: %w", err)
	}
	return lo.Map(
		users,
		func(item user, _ int) models.User {
			return item.toServiceModel()
		},
	), nil
}

const usersRepoGetListByRoleIdCountQuery = `
select count(*)
from public.user
where role_id = $1
`

func (r *UsersRepo) GetListByRoleIdCount(
	ctx context.Context,
	opts repo.UsersRepoGetListByRoleIdOpts,
) (int, error) {
	var count int
	if err := r.db.GetContext(ctx, &count, usersRepoGetListByRoleIdCountQuery, opts.RoleId); err != nil {
		return 0, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return count, nil
}

const getByCredentialsQuery = `
select 
    u.id, 
    u.group_id, 
    u.role_id,
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
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repo.ErrNotFound
		}
		return models.User{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return u.toServiceModel(), nil
}

const createQuery = `
insert into public.user (group_id, role_id, email, password, first_name, last_name, middle_name) 
values (:group_id, :role_id, :email, :password, :first_name, :last_name, :middle_name)
`

func (r *UsersRepo) Create(
	ctx context.Context,
	opts repo.UsersRepoCreateOpts,
) (models.User, error) {
	_, err := r.db.NamedExecContext(ctx, createQuery, struct {
		GroupId    *int64  `db:"group_id"`
		RoleId     int64   `db:"role_id"`
		Email      string  `db:"email"`
		Password   string  `db:"password"`
		FirstName  string  `db:"first_name"`
		LastName   string  `db:"last_name"`
		MiddleName *string `db:"middle_name"`
	}{
		GroupId:    opts.GroupId,
		RoleId:     opts.RoleId,
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
