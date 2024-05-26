package pg

import (
	"backend/internal/repo"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type TaskLinksRepo struct {
	db *sqlx.DB
}

func NewTaskLinksRepo(db *sqlx.DB) *TaskLinksRepo {
	return &TaskLinksRepo{db: db}
}

const taskLinksCreateQuery = `
insert into task_links (user_id, group_id, task_id) values ($1, $2, $3)
`

func (r *TaskLinksRepo) Create(
	ctx context.Context,
	opts repo.TaskLinksRepoCreateOpts,
) error {
	if _, err := r.db.ExecContext(ctx, taskLinksCreateQuery, opts.UserId, opts.GroupId, opts.TaskId); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}
	return nil
}
