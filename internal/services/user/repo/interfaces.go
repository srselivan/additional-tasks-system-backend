package repo

import (
	"backend/internal/models"
	"context"
)

type UsersRepo interface {
	GetById(
		ctx context.Context,
		id int64,
	) (models.User, error)
	GetByCredentials(
		ctx context.Context,
		opts UsersRepoGetByCredentialsOpts,
	) (models.User, error)
	Create(
		ctx context.Context,
		opts UsersRepoCreateOpts,
	) (models.User, error)
	//Update(
	//	ctx context.Context,
	//	opts UsersRepoUpdateOpts,
	//) (models.User, error)
}
