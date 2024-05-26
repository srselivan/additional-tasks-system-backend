package pg

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type task struct {
	Id            int64      `db:"id"`
	Title         string     `db:"title"`
	Text          *string    `db:"text"`
	CreatedBy     int64      `db:"created_by"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	EffectiveFrom time.Time  `db:"effective_from"`
	EffectiveTill time.Time  `db:"effective_till"`
}

func (t task) toServiceModel() models.Task {
	return models.Task{
		Id:            t.Id,
		Title:         t.Title,
		Text:          t.Text,
		CreatedBy:     t.CreatedBy,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
		EffectiveFrom: t.EffectiveFrom,
		EffectiveTill: t.EffectiveTill,
	}
}

type TasksRepo struct {
	db *sqlx.DB
}

func NewTasksRepo(db *sqlx.DB) *TasksRepo {
	return &TasksRepo{db: db}
}

const tasksRepoGetByIdQuery = `
select 
    t.id,
    t.title, 
    t.text, 
    t.created_by, 
    t.effective_from, 
    t.effective_till, 
    t.created_at, 
    t.updated_at
from public.task t
where t.id = $1
`

func (r *TasksRepo) GetById(
	ctx context.Context,
	id int64,
) (models.Task, error) {
	var t task
	if err := r.db.GetContext(ctx, &t, tasksRepoGetByIdQuery, id); err != nil {
		return models.Task{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return t.toServiceModel(), nil
}

const tasksRepoGetListForCreatorQuery = `
select 
    t.id, 
    t.created_by, 
    t.title, 
    t.text, 
    t.effective_from, 
    t.effective_till, 
    t.created_at, 
    t.updated_at
from public.task t
where t.created_by = $1
order by id desc 
limit $2
offset $3
`

func (r *TasksRepo) GetListForCreator(
	ctx context.Context,
	opts repo.TasksRepoGetListForCreatorOpts,
) ([]models.Task, error) {
	var tasks []task
	if err := r.db.SelectContext(ctx, &tasks, tasksRepoGetListForCreatorQuery, opts.CreatedBy, opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("r.db.SelectContext: %w", err)
	}
	return lo.Map(
		tasks,
		func(item task, _ int) models.Task {
			return item.toServiceModel()
		},
	), nil
}

const tasksRepoGetCountForCreatorQuery = `
select count(*)
from public.task t
where t.created_by = $1
`

func (r *TasksRepo) GetCountForCreator(
	ctx context.Context,
	createdBy int64,
) (int64, error) {
	var count int64
	if err := r.db.GetContext(ctx, &count, tasksRepoGetCountForCreatorQuery, createdBy); err != nil {
		return 0, fmt.Errorf("r.db.GetContext: :%w", err)
	}
	return count, nil
}

const tasksRepoGetListForUserQuery = `
select 
    t.id, 
    t.created_by, 
    t.title, 
    t.text, 
    t.effective_from, 
    t.effective_till, 
    t.created_at, 
    t.updated_at
from public.task t
left join public.task_links tl on t.id = tl.task_id
where tl.user_id = $1 or tl.group_id = $2
order by id desc 
limit $3
offset $4
`

func (r *TasksRepo) GetListForUser(
	ctx context.Context,
	opts repo.TasksRepoGetListForUserOpts,
) ([]models.Task, error) {
	var tasks []task
	if err := r.db.SelectContext(
		ctx, &tasks, tasksRepoGetListForUserQuery,
		opts.UserId,
		opts.GroupId,
		opts.Limit,
		opts.Offset,
	); err != nil {
		return nil, fmt.Errorf("r.db.SelectContext: %w", err)
	}
	return lo.Map(
		tasks,
		func(item task, _ int) models.Task {
			return item.toServiceModel()
		},
	), nil
}

const tasksRepoGetCountForUserQuery = `
select count(*)
from public.task t
left join public.task_links tl on t.id = tl.task_id
where tl.user_id = $1 or tl.group_id = $2
`

func (r *TasksRepo) GetCountForUser(
	ctx context.Context,
	userId, groupId int64,
) (int64, error) {
	var count int64
	if err := r.db.GetContext(ctx, &count, tasksRepoGetCountForUserQuery, userId, groupId); err != nil {
		return 0, fmt.Errorf("r.db.GetContext: :%w", err)
	}
	return count, nil
}

const tasksRepoCreateQuery = `
insert into public.task (created_by, title, text, effective_from, effective_till) 
values (:created_by, :title, :text, :effective_from, :effective_till)
returning id, created_at, title, text, effective_from, effective_till, created_at, updated_at
`

func (r *TasksRepo) Create(
	ctx context.Context,
	opts repo.TasksRepoCreateOpts,
) (models.Task, error) {
	rows, err := r.db.NamedQueryContext(ctx, tasksRepoCreateQuery, struct {
		CreatedBy     int64     `db:"created_by"`
		Title         string    `db:"title"`
		Text          *string   `db:"text"`
		EffectiveFrom time.Time `db:"effective_from"`
		EffectiveTill time.Time `db:"effective_till"`
	}{
		CreatedBy:     opts.CreatedBy,
		Title:         opts.Title,
		Text:          opts.Text,
		EffectiveFrom: opts.EffectiveFrom,
		EffectiveTill: opts.EffectiveTill,
	})
	if err != nil {
		return models.Task{}, fmt.Errorf("r.db.NamedQueryContext: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return models.Task{}, errors.New("no rows returned")
	}
	var t task
	if err = rows.StructScan(&t); err != nil {
		return models.Task{}, fmt.Errorf("rows.StructScan: %w", err)
	}
	if rows.Err() != nil {
		return models.Task{}, fmt.Errorf("rows.Err: %w", err)
	}

	return t.toServiceModel(), nil
}

const tasksRepoUpdateQuery = `
update public.task
set (
     title, 
     text, 
     effective_from, 
     effective_till
    ) = (
     :title, 
     :text, 
     :effective_from, 
     :effective_till
    )
where id = :id
`

func (r *TasksRepo) Update(
	ctx context.Context,
	opts repo.TasksRepoUpdateOpts,
) (models.Task, error) {
	_, err := r.db.NamedExecContext(ctx, tasksRepoUpdateQuery, struct {
		Id            int64     `db:"id"`
		GroupId       int64     `db:"group_id"`
		Title         string    `db:"title"`
		Text          string    `db:"text"`
		EffectiveFrom time.Time `db:"effective_from"`
		EffectiveTill time.Time `db:"effective_till"`
		Cost          int64     `db:"cost"`
	}{
		Id:            opts.Id,
		GroupId:       opts.GroupId,
		Title:         opts.Title,
		Text:          opts.Text,
		EffectiveFrom: opts.EffectiveFrom,
		EffectiveTill: opts.EffectiveTill,
		Cost:          opts.Cost,
	})
	if err != nil {
		return models.Task{}, fmt.Errorf("r.db.NamedExecContext: %w", err)
	}

	return models.Task{}, nil
}

const tasksRepoDeleteQuery = `
delete from public.task where id = $1
`

func (r *TasksRepo) Delete(
	ctx context.Context,
	id int64,
) error {
	if _, err := r.db.ExecContext(ctx, tasksRepoDeleteQuery, id); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}
	return nil
}
