package repo

import (
	"backend/internal/models"
	"context"
)

type TasksRepo interface {
	GetById(
		ctx context.Context,
		id int64,
	) (models.Task, error)
	GetList(
		ctx context.Context,
		opts TasksRepoGetListOpts,
	) ([]models.Task, error)
	GetCount(
		ctx context.Context,
		groupId int64,
	) (int64, error)
	Create(
		ctx context.Context,
		opts TasksRepoCreateOpts,
	) (models.Task, error)
	Update(
		ctx context.Context,
		opts TasksRepoUpdateOpts,
	) (models.Task, error)
	Delete(
		ctx context.Context,
		id int64,
	) error
}
