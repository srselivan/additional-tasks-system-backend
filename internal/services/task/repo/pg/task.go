package pg

import (
	"backend/internal/models"
	"backend/internal/services/task/repo"
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type TasksRepo struct {
	db *sqlx.DB
}

func NewTasksRepo(db *sqlx.DB) *TasksRepo {
	return &TasksRepo{db: db}
}

const getByIdQuery = `
select 
    t.id, 
    t.group_id, 
    t.title, 
    t.text, 
    t.cost, 
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
	if err := r.db.GetContext(ctx, &t, getByIdQuery, id); err != nil {
		return models.Task{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return t.toServiceModel(), nil
}

const getListQuery = `
select 
    t.id, 
    t.group_id, 
    t.title, 
    t.text, 
    t.cost, 
    t.effective_from, 
    t.effective_till, 
    t.created_at, 
    t.updated_at
from public.task t
where t.group_id = $1
order by id desc 
limit $2
offset $3
`

func (r *TasksRepo) GetList(
	ctx context.Context,
	opts repo.TasksRepoGetListOpts,
) ([]models.Task, error) {
	var tasks []task
	if err := r.db.SelectContext(ctx, &tasks, getListQuery, opts.GroupId, opts.Limit, opts.Offset); err != nil {
		return nil, fmt.Errorf("r.db.SelectContext: %w", err)
	}
	return lo.Map(
		tasks,
		func(item task, _ int) models.Task {
			return item.toServiceModel()
		},
	), nil
}

const getCountQuery = `
select count(*)
from public.task t
where t.group_id = $1
`

func (r *TasksRepo) GetCount(
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
insert into public.task (group_id, title, text, cost, effective_from, effective_till) 
values (:group_id, :title, :text, :cost, :effective_from, :effective_till)
`

func (r *TasksRepo) Create(
	ctx context.Context,
	opts repo.TasksRepoCreateOpts,
) (models.Task, error) {
	_, err := r.db.NamedExecContext(ctx, createQuery, struct {
		GroupId       int64     `db:"group_id"`
		Title         string    `db:"title"`
		Text          string    `db:"text"`
		EffectiveFrom time.Time `db:"effective_from"`
		EffectiveTill time.Time `db:"effective_till"`
		Cost          int64     `db:"cost"`
	}{
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

const updateQuery = `
update public.task
set (
     group_id, 
     title, 
     text, 
     effective_from, 
     effective_till, 
     cost
    ) = (
     :group_id, 
     :title, 
     :text, 
     :effective_from, 
     :effective_till, 
     :cost
    )
where id = :id
`

func (r *TasksRepo) Update(
	ctx context.Context,
	opts repo.TasksRepoUpdateOpts,
) (models.Task, error) {
	_, err := r.db.NamedExecContext(ctx, updateQuery, struct {
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

const deleteQuery = `
delete from public.task where id = $1
`

func (r *TasksRepo) Delete(
	ctx context.Context,
	id int64,
) error {
	if _, err := r.db.ExecContext(ctx, deleteQuery, id); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}
	return nil
}
