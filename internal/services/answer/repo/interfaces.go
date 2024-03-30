package repo

import (
	"backend/internal/models"
	"context"
)

type AnswersRepo interface {
	GetById(
		ctx context.Context,
		id int64,
	) (models.Answer, error)
	GetList(
		ctx context.Context,
		opts AnswersRepoGetListOpts,
	) ([]models.Answer, error)
	GetCount(
		ctx context.Context,
		groupId int64,
	) (int64, error)
	Create(
		ctx context.Context,
		opts AnswersRepoCreateOpts,
	) (models.Answer, error)
	Update(
		ctx context.Context,
		opts AnswersRepoUpdateOpts,
	) (models.Answer, error)
	Delete(
		ctx context.Context,
		id int64,
	) error
}
