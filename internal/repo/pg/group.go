package pg

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type group struct {
	Id        int64      `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (g group) toServiceModel() models.Group {
	return models.Group{
		Id:        g.Id,
		Name:      g.Name,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}

type GroupsRepo struct {
	db *sqlx.DB
}

func NewGroupsRepo(db *sqlx.DB) *GroupsRepo {
	return &GroupsRepo{db: db}
}

const groupsRepoGetByIdQuery = `
select 
   g.id, 
   g.name, 
   g.created_at, 
   g.updated_at
from public."group" g
where g.id = $1
`

func (r *GroupsRepo) GetById(
	ctx context.Context,
	id int64,
) (models.Group, error) {
	var g group
	if err := r.db.GetContext(ctx, &g, groupsRepoGetByIdQuery, id); err != nil {
		return models.Group{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return g.toServiceModel(), nil
}

const groupsRepoGetListQuery = `
select 
   g.id, 
   g.name, 
   g.created_at, 
   g.updated_at
from public."group" g
order by id desc 
limit $1
offset $2
`

func (r *GroupsRepo) GetList(
	ctx context.Context,
	opts repo.GroupsRepoGetListOpts,
) ([]models.Group, error) {
	var groups []group
	if err := r.db.SelectContext(ctx, &groups, groupsRepoGetListQuery, opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("r.db.SelectContext: %w", err)
	}
	return lo.Map(
		groups,
		func(item group, _ int) models.Group {
			return item.toServiceModel()
		},
	), nil
}

const groupsRepoGetCountQuery = `
select count(*)
from public."group"
`

func (r *GroupsRepo) GetCount(
	ctx context.Context,
) (int64, error) {
	var count int64
	if err := r.db.GetContext(ctx, &count, groupsRepoGetCountQuery); err != nil {
		return 0, fmt.Errorf("r.db.GetContext: :%w", err)
	}
	return count, nil
}

const groupsRepoCreateQuery = `
insert into public."group" (name)
values (:name)
`

func (r *GroupsRepo) Create(
	ctx context.Context,
	opts repo.GroupsRepoCreateOpts,
) (models.Group, error) {
	_, err := r.db.NamedExecContext(ctx, groupsRepoCreateQuery, struct {
		Name string `db:"name"`
	}{
		Name: opts.Name,
	})
	if err != nil {
		return models.Group{}, fmt.Errorf("r.db.NamedExecContext: %w", err)
	}

	return models.Group{}, nil
}

const groupsRepoUpdateQuery = `
update public."group"
set (
     name
    ) = (
     :name
    )
where id = :id
`

func (r *GroupsRepo) Update(
	ctx context.Context,
	opts repo.GroupsRepoUpdateOpts,
) (models.Group, error) {
	_, err := r.db.NamedExecContext(ctx, groupsRepoUpdateQuery, struct {
		Id   int64  `db:"id"`
		Name string `db:"name"`
	}{
		Id:   opts.Id,
		Name: opts.Name,
	})
	if err != nil {
		return models.Group{}, fmt.Errorf("r.db.NamedExecContext: %w", err)
	}

	return models.Group{}, nil
}

const groupsRepoDeleteQuery = `
delete from public."group" where id = $1
`

func (r *GroupsRepo) Delete(
	ctx context.Context,
	id int64,
) error {
	if _, err := r.db.ExecContext(ctx, groupsRepoDeleteQuery, id); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}
	return nil
}
