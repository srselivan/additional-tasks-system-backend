package pg

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type answer struct {
	Id        int64      `db:"id"`
	UserId    int64      `db:"user_id"`
	TaskId    int64      `db:"task_id"`
	Comment   *string    `db:"comment"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (a answer) toServiceModel() models.Answer {
	return models.Answer{
		Id:        a.Id,
		UserId:    a.UserId,
		TaskId:    a.TaskId,
		Comment:   a.Comment,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

type AnswersRepo struct {
	db *sqlx.DB
}

func NewAnswersRepo(db *sqlx.DB) *AnswersRepo {
	return &AnswersRepo{db: db}
}

const answersRepGetByIdQuery = `
select 
    a.id, 
    a.user_id,
    a.task_id,
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
	if err := r.db.GetContext(ctx, &a, answersRepGetByIdQuery, id); err != nil {
		return models.Answer{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return a.toServiceModel(), nil
}

const answersRepGetByTaskIdQuery = `
select 
    a.id, 
    a.user_id,
    a.task_id,
    a.comment, 
    a.created_at, 
    a.updated_at
from public.answer a
where a.task_id = $1 and a.user_id = $2
`

func (r *AnswersRepo) GetByTaskIdAndUserId(
	ctx context.Context,
	userId int64,
	taskId int64,
) (models.Answer, error) {
	var a answer
	if err := r.db.GetContext(ctx, &a, answersRepGetByTaskIdQuery, taskId, userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Answer{}, repo.ErrNotFound
		}
		return models.Answer{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return a.toServiceModel(), nil
}

const answersRepoGetListQuery = `
select 
    a.id, 
    a.task_id,
    a.user_id,
    a.comment, 
    a.created_at, 
    a.updated_at
from public.answer a
where task_id = $1
order by id desc 
limit $2
offset $3
`

func (r *AnswersRepo) GetList(
	ctx context.Context,
	opts repo.AnswersRepoGetListOpts,
) ([]models.Answer, error) {
	var answers []answer
	if err := r.db.SelectContext(ctx, &answers, answersRepoGetListQuery, opts.TaskId, opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("r.db.SelectContext: %w", err)
	}
	return lo.Map(
		answers,
		func(item answer, _ int) models.Answer {
			return item.toServiceModel()
		},
	), nil
}

const answersRepoGetCountQuery = `
select count(*)
from public.answer a
where a.task_id = $1
`

func (r *AnswersRepo) GetCount(
	ctx context.Context,
	groupId int64,
) (int64, error) {
	var count int64
	if err := r.db.GetContext(ctx, &count, answersRepoGetCountQuery, groupId); err != nil {
		return 0, fmt.Errorf("r.db.GetContext: :%w", err)
	}
	return count, nil
}

const answersRepoCreateQuery = `
insert into public.answer (task_id, user_id, comment)
values (:task_id, :user_id, :comment)
returning id
`

func (r *AnswersRepo) Create(
	ctx context.Context,
	opts repo.AnswersRepoCreateOpts,
) (models.Answer, error) {
	rows, err := r.db.NamedQueryContext(ctx, answersRepoCreateQuery, struct {
		TaskId  int64   `db:"task_id"`
		UserId  int64   `db:"user_id"`
		Comment *string `db:"comment"`
	}{
		TaskId:  opts.TaskId,
		UserId:  opts.UserId,
		Comment: opts.Comment,
	})
	if err != nil {
		return models.Answer{}, fmt.Errorf("r.db.NamedExecContext: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return models.Answer{}, errors.New("empty result")
	}

	var id int64
	if err = rows.Scan(&id); err != nil {
		return models.Answer{}, fmt.Errorf("rows.Scan: %w", err)
	}

	if err = rows.Err(); err != nil {
		return models.Answer{}, fmt.Errorf("rows.Err: %w", err)
	}

	return models.Answer{
		Id:        id,
		UserId:    0,
		TaskId:    0,
		Comment:   nil,
		CreatedAt: time.Time{},
		UpdatedAt: nil,
	}, nil
}

const answersRepoUpdateQuery = `
update public.answer
set comment = :comment
where id = :id
`

func (r *AnswersRepo) Update(
	ctx context.Context,
	opts repo.AnswersRepoUpdateOpts,
) (models.Answer, error) {
	_, err := r.db.NamedExecContext(ctx, answersRepoUpdateQuery, struct {
		Id      int64  `db:"id"`
		Comment string `db:"comment"`
	}{
		Id:      opts.Id,
		Comment: opts.Comment,
	})
	if err != nil {
		return models.Answer{}, fmt.Errorf("r.db.NamedExecContext: %w", err)
	}

	return models.Answer{}, nil
}

const answersRepoDeleteQuery = `
delete from public.answer where id = $1
`

func (r *AnswersRepo) Delete(
	ctx context.Context,
	id int64,
) error {
	if _, err := r.db.ExecContext(ctx, answersRepoDeleteQuery, id); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}
	return nil
}
