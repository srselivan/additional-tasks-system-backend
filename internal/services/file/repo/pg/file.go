package pg

import (
	"backend/internal/models"
	"backend/internal/services/file/repo"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type FilesRepo struct {
	db *sqlx.DB
}

func NewFilesRepo(db *sqlx.DB) *FilesRepo {
	return &FilesRepo{db: db}
}

const getByIdQuery = `
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
	if err := r.db.GetContext(ctx, &f, getByIdQuery, id); err != nil {
		return models.File{}, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return f.toServiceModel(), nil
}

const createQuery = `
insert into public.file (name, filename, filepath, task_id, answer_id) 
values (:name, :filename, :filepath, :task_id, :answer_id)
`

func (r *FilesRepo) Create(
	ctx context.Context,
	opts repo.FilesRepoCreateOpts,
) (models.File, error) {
	_, err := r.db.NamedExecContext(ctx, createQuery, struct {
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

	return models.File{}, nil
}

const deleteQuery = `
delete from public.file where id = $1
`

func (r *FilesRepo) Delete(
	ctx context.Context,
	id int64,
) error {
	if _, err := r.db.ExecContext(ctx, deleteQuery, id); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}
	return nil
}
