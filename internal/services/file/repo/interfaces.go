package repo

import (
	"backend/internal/models"
	"context"
)

type FilesRepo interface {
	GetById(
		ctx context.Context,
		id int64,
	) (models.File, error)
	Create(
		ctx context.Context,
		opts FilesRepoCreateOpts,
	) (models.File, error)
	Delete(
		ctx context.Context,
		id int64,
	) error
}
