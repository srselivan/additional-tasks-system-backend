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

type file struct {
	Id        int64      `db:"id"`
	Name      string     `db:"name"`
	Filename  string     `db:"filename"`
	Filepath  string     `db:"filepath"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	TaskId    *int64     `db:"task_id"`
	AnswerId  *int64     `db:"answer_id"`
}

func (f file) toServiceModel() models.File {
	return models.File{
		Id:        f.Id,
		Name:      f.Name,
		Filename:  f.Filename,
		Filepath:  f.Filepath,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		TaskId:    f.TaskId,
		AnswerId:  f.AnswerId,
	}
}

type FilesRepo struct {
	db *sqlx.DB
}

func NewFilesRepo(db *sqlx.DB) *FilesRepo {
	return &FilesRepo{db: db}
}

const filesRepoGetByIdQuery = `
select 
    f.id, 
    f.name, 
    f.filename, 
    f.filepath, 
    f.created_at, 
    f.updated_at, 
    f.task_id, 
    f.answer_id
from public.file f
where f.id = $1
`

func (r *FilesRepo) GetById(
	ctx context.Context,
	id int64,
) (models.File, error) {
	var f file
	if err := r.db.GetContext(ctx, &f, filesRepoGetByIdQuery, id); err != nil {
		return models.File{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return f.toServiceModel(), nil
}

const filesRepoGetByAnswerIdQuery = `
select 
    f.id, 
    f.name, 
    f.filename, 
    f.filepath, 
    f.created_at, 
    f.updated_at, 
    f.task_id, 
    f.answer_id
from public.file f
where f.answer_id = $1
`

func (r *FilesRepo) GetByAnswerId(
	ctx context.Context,
	answerId int64,
) ([]models.File, error) {
	var files []file
	if err := r.db.SelectContext(ctx, &files, filesRepoGetByAnswerIdQuery, answerId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrNotFound
		}
		return nil, fmt.Errorf("r.db.SelectContext: %w", err)
	}
	return lo.Map(
		files,
		func(item file, _ int) models.File {
			return item.toServiceModel()
		},
	), nil
}

const filesRepoGetByTaskIdQuery = `
select 
    f.id, 
    f.name, 
    f.filename, 
    f.filepath, 
    f.created_at, 
    f.updated_at, 
    f.task_id, 
    f.answer_id
from public.file f
where f.task_id = $1
`

func (r *FilesRepo) GetByTaskId(
	ctx context.Context,
	taskId int64,
) ([]models.File, error) {
	var files []file
	if err := r.db.SelectContext(ctx, &files, filesRepoGetByTaskIdQuery, taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrNotFound
		}
		return nil, fmt.Errorf("r.db.SelectContext: %w", err)
	}
	return lo.Map(
		files,
		func(item file, _ int) models.File {
			return item.toServiceModel()
		},
	), nil
}

const filesRepoCreateQuery = `
insert into public.file (name, filename, filepath, task_id, answer_id) 
values (:name, :filename, :filepath, :task_id, :answer_id)
returning id
`

func (r *FilesRepo) Create(
	ctx context.Context,
	opts repo.FilesRepoCreateOpts,
) (models.File, error) {
	rows, err := r.db.NamedQueryContext(ctx, filesRepoCreateQuery, struct {
		Name     string `db:"name"`
		Filename string `db:"filename"`
		Filepath string `db:"filepath"`
		TaskId   *int64 `db:"task_id"`
		AnswerId *int64 `db:"answer_id"`
	}{
		Name:     opts.Name,
		Filename: opts.Filename,
		Filepath: opts.Filepath,
		TaskId:   opts.TaskId,
		AnswerId: opts.AnswerId,
	})
	if err != nil {
		return models.File{}, fmt.Errorf("r.db.NamedExecContext: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return models.File{}, errors.New("empty result")
	}

	var id int64
	if err = rows.Scan(&id); err != nil {
		return models.File{}, fmt.Errorf("rows.Scan: %w", err)
	}

	if err = rows.Err(); err != nil {
		return models.File{}, fmt.Errorf("rows.Err: %w", err)
	}

	return models.File{
		Id:        id,
		Name:      opts.Name,
		Filename:  opts.Filename,
		Filepath:  opts.Filepath,
		CreatedAt: time.Time{},
		UpdatedAt: nil,
		TaskId:    opts.TaskId,
		AnswerId:  nil,
	}, nil
}

const filesRepoUpdateAnswerIdQuery = `
update public.file set answer_id = $1 where id = $2
`

func (r *FilesRepo) UpdateAnswerId(
	ctx context.Context,
	answerId int64,
	fileId int64,
) error {
	if _, err := r.db.ExecContext(ctx, filesRepoUpdateAnswerIdQuery, answerId, fileId); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}
	return nil
}

const filesRepoUpdateTaskAnswerIdQuery = `
update public.file set task_id = $1 where id = $2
`

func (r *FilesRepo) UpdateTaskId(
	ctx context.Context,
	taskId int64,
	fileId int64,
) error {
	if _, err := r.db.ExecContext(ctx, filesRepoUpdateTaskAnswerIdQuery, taskId, fileId); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}
	return nil
}

const filesRepoDeleteQuery = `
delete from public.file where id = $1
`

func (r *FilesRepo) Delete(
	ctx context.Context,
	id int64,
) error {
	if _, err := r.db.ExecContext(ctx, filesRepoDeleteQuery, id); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}
	return nil
}
