package repo

import (
	"backend/internal/models"
	"context"
)

type AnswersRepo interface {
	GetById(
		ctx context.Context,
		id int64,
	) (models.Task, error)
	GetList(
		ctx context.Context,
		opts AnswersRepoGetListOpts,
	) ([]models.Task, error)
	GetCount(
		ctx context.Context,
		groupId int64,
	) (int64, error)
	Create(
		ctx context.Context,
		opts AnswersRepoCreateOpts,
	) (models.Task, error)
	Update(
		ctx context.Context,
		opts AnswersRepoUpdateOpts,
	) (models.Task, error)
	Delete(
		ctx context.Context,
		id int64,
	) error
}
