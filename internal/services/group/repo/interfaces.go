package repo

import (
	"backend/internal/models"
	"context"
)

type GroupsRepo interface {
	GetById(
		ctx context.Context,
		id int64,
	) (models.Group, error)
	GetList(
		ctx context.Context,
		opts GroupsRepoGetListOpts,
	) ([]models.Group, error)
	GetCount(
		ctx context.Context,
	) (int64, error)
	Create(
		ctx context.Context,
		opts GroupsRepoCreateOpts,
	) (models.Group, error)
	Update(
		ctx context.Context,
		opts GroupsRepoUpdateOpts,
	) (models.Group, error)
	Delete(
		ctx context.Context,
		id int64,
	) error
}
