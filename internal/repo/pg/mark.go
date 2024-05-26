package pg

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"time"
)

type mark struct {
	Id        int64      `db:"id"`
	AnswerId  int64      `db:"answer_id"`
	TaskId    int64      `db:"task_id"`
	Mark      int64      `db:"mark"`
	Comment   *string    `db:"comment"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (m mark) toServiceModel() models.Mark {
	return models.Mark{
		Id:        m.Id,
		AnswerId:  m.AnswerId,
		TaskId:    m.TaskId,
		Mark:      m.Mark,
		Comment:   m.Comment,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

type MarksRepo struct {
	db *sqlx.DB
}

func NewMarksRepo(db *sqlx.DB) *MarksRepo {
	return &MarksRepo{db: db}
}

const marksRepoGetByIdQuery = `
select m.id, m.answer_id, m.mark, m.comment, m.created_at, m.updated_at, a.task_id
from public.mark m
left join public.answer a on a.id = m.answer_id
where m.id = $1
`

func (r *MarksRepo) GetById(
	ctx context.Context,
	id int64,
) (models.Mark, error) {
	var m mark
	if err := r.db.GetContext(ctx, &m, marksRepoGetByIdQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Mark{}, repo.ErrNotFound
		}
		return models.Mark{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return m.toServiceModel(), nil
}

const marksRepoGetByAnswerIdQuery = `
select m.id, m.answer_id, m.mark, m.comment, m.created_at, m.updated_at, a.task_id
from public.mark m
left join public.answer a on a.id = m.answer_id
where m.answer_id = $1
`

func (r *MarksRepo) GetByAnswerId(
	ctx context.Context,
	id int64,
) (models.Mark, error) {
	var m mark
	if err := r.db.GetContext(ctx, &m, marksRepoGetByAnswerIdQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Mark{}, repo.ErrNotFound
		}
		return models.Mark{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return m.toServiceModel(), nil
}

const marksRepoGetListByUserIdQuery = `
select m.id, m.answer_id, m.mark, m.comment, m.created_at, m.updated_at, a.task_id
from public.mark m
left join public.answer a on a.id = m.answer_id
where a.user_id = $1
order by m.id
limit $2
offset $3
`

func (r *MarksRepo) GetListByUserId(
	ctx context.Context,
	opts repo.MarksRepoGetListByUserIdOpts,
) ([]models.Mark, error) {
	var ms []mark
	if err := r.db.SelectContext(ctx, &ms, marksRepoGetListByUserIdQuery, opts.UserId, opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("r.db.SelectContext: %w", err)
	}
	return lo.Map(
		ms,
		func(item mark, _ int) models.Mark {
			return item.toServiceModel()
		},
	), nil
}

const marksRepoGetCountByUserIdQuery = `
select count(*)
from public.mark m 
left join public.answer a on a.id = m.answer_id
where a.user_id = $1
`

func (r *MarksRepo) GetCountByUserId(
	ctx context.Context,
	userId int64,
) (int64, error) {
	var count int64
	if err := r.db.GetContext(ctx, &count, marksRepoGetCountByUserIdQuery, userId); err != nil {
		return 0, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return count, nil
}

const marksRepoCreateQuery = `
insert into public.mark (mark, comment, answer_id)
values ($1, $2, $3)
returning id, mark, comment, answer_id, created_at, updated_at
`

func (r *MarksRepo) Create(
	ctx context.Context,
	opts repo.MarksRepoCreateOpts,
) (models.Mark, error) {
	rows, err := r.db.QueryxContext(ctx, marksRepoCreateQuery, opts.Mark, opts.Comment, opts.AnswerId)
	if err != nil {
		return models.Mark{}, fmt.Errorf("r.db.QueryContext: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return models.Mark{}, errors.New("no rows returned")
	}
	var m mark
	if err = rows.StructScan(&m); err != nil {
		return models.Mark{}, fmt.Errorf("rows.StructScan: %w", err)
	}
	if rows.Err() != nil {
		return models.Mark{}, fmt.Errorf("rows.Err: %w", err)
	}

	return m.toServiceModel(), nil
}

const marksRepoUpdateQuery = `
update public.mark set mark = $1, comment = $2, updated_at = now()
where id = $3
returning id, mark, comment, answer_id, created_at, updated_at
`

func (r *MarksRepo) Update(
	ctx context.Context,
	opts repo.MarksRepoUpdateOpts,
) (models.Mark, error) {
	rows, err := r.db.QueryxContext(ctx, marksRepoUpdateQuery, opts.Mark, opts.Comment, opts.Id)
	if err != nil {
		return models.Mark{}, fmt.Errorf("r.db.QueryxContext: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return models.Mark{}, errors.New("no rows returned")
	}
	var m mark
	if err = rows.StructScan(&m); err != nil {
		return models.Mark{}, fmt.Errorf("rows.StructScan: %w", err)
	}
	if rows.Err() != nil {
		return models.Mark{}, fmt.Errorf("rows.Err: %w", err)
	}

	return m.toServiceModel(), nil
}

const marksRepoDeleteQuery = `
delete from public.mark where id = $1
`

func (r *MarksRepo) Delete(
	ctx context.Context,
	id int64,
) error {
	if _, err := r.db.ExecContext(ctx, marksRepoDeleteQuery, id); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}
	return nil
}
