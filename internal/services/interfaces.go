package services

import (
	"backend/internal/models"
	"context"
)

type TaskService interface {
	GetById(
		ctx context.Context,
		id int64,
	) (models.Task, error)
	GetList(
		ctx context.Context,
		opts TaskServiceGetListOpts,
	) ([]models.Task, error)
	GetCount(
		ctx context.Context,
		groupId int64,
	) (int64, error)
	Create(
		ctx context.Context,
		opts TaskServiceCreateOpts,
	) (models.Task, error)
	Update(
		ctx context.Context,
		opts TaskServiceUpdateOpts,
	) (models.Task, error)
	Delete(
		ctx context.Context,
		id int64,
	) error
}
