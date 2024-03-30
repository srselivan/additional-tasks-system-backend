package pg

import (
	"backend/internal/models"
	"backend/internal/services/answer/repo"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type AnswersRepo struct {
	db *sqlx.DB
}

func NewAnswersRepo(db *sqlx.DB) *AnswersRepo {
	return &AnswersRepo{db: db}
}

const getByIdQuery = `
select 
    a.id, 
    a.group_id, 
    a.comment, 
    a.created_at, 
    a.updated_at
from public.answer a
where a.id = $1
`

func (r *AnswersRepo) GetById(
	ctx context.Context,
	id int64,
) (models.Answer, error) {
	var a answer
	if err := r.db.GetContext(ctx, &a, getByIdQuery, id); err != nil {
		return models.Answer{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return a.toServiceModel(), nil
}

const getListQuery = `
select 
    a.id, 
    a.group_id, 
    a.comment, 
    a.created_at, 
    a.updated_at
from public.answer a
where a.group_id = $1
order by id desc 
limit $2
offset $3
`

func (r *AnswersRepo) GetList(
	ctx context.Context,
	opts repo.AnswersRepoGetListOpts,
) ([]models.Answer, error) {
	var answers []answer
	if err := r.db.SelectContext(ctx, &answers, getListQuery, opts.GroupId, opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("r.db.SelectContext: %w", err)
	}
	return lo.Map(
		answers,
		func(item answer, _ int) models.Answer {
			return item.toServiceModel()
		},
	), nil
}

const getCountQuery = `
select count(*)
from public.answer a
where a.group_id = $1
`

func (r *AnswersRepo) GetCount(
	ctx context.Context,
	groupId int64,
) (int64, error) {
	var count int64
	if err := r.db.GetContext(ctx, &count, getCountQuery, groupId); err != nil {
		return 0, fmt.Errorf("r.db.GetContext: :%w", err)
	}
	return count, nil
}

const createQuery = `
insert into public.answer (group_id, comment)
values (group_id, comment)
`

func (r *AnswersRepo) Create(
	ctx context.Context,
	opts repo.AnswersRepoCreateOpts,
) (models.Answer, error) {
	_, err := r.db.NamedExecContext(ctx, createQuery, struct {
		GroupId int64  `db:"group_id"`
		Comment string `db:"comment"`
	}{
		GroupId: opts.GroupId,
		Comment: opts.Comment,
	})
	if err != nil {
		return models.Answer{}, fmt.Errorf("r.db.NamedExecContext: %w", err)
	}

	return models.Answer{}, nil
}

const updateQuery = `
update public.answer
set (
     group_id,
     comment
    ) = (
     :group_id, 
     :comment
    )
where id = :id
`

func (r *AnswersRepo) Update(
	ctx context.Context,
	opts repo.AnswersRepoUpdateOpts,
) (models.Answer, error) {
	_, err := r.db.NamedExecContext(ctx, updateQuery, struct {
		Id      int64  `db:"id"`
		GroupId int64  `db:"group_id"`
		Comment string `db:"comment"`
	}{
		Id:      opts.Id,
		GroupId: opts.GroupId,
		Comment: opts.Comment,
	})
	if err != nil {
		return models.Answer{}, fmt.Errorf("r.db.NamedExecContext: %w", err)
	}

	return models.Answer{}, nil
}

const deleteQuery = `
delete from public.answer where id = $1
`

func (r *AnswersRepo) Delete(
	ctx context.Context,
	id int64,
) error {
	if _, err := r.db.ExecContext(ctx, deleteQuery, id); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}
	return nil
}
