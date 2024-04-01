package services

import (
	"backend/internal/models"
	"context"
)

type TaskService interface {
	GetById(ctx context.Context, id int64) (models.Task, error)
	GetList(ctx context.Context, opts TaskServiceGetListOpts) ([]models.Task, error)
	GetCount(ctx context.Context, groupId int64) (int64, error)
	Create(ctx context.Context, opts TaskServiceCreateOpts) (models.Task, error)
	Update(ctx context.Context, opts TaskServiceUpdateOpts) (models.Task, error)
	Delete(ctx context.Context, id int64) error
}

type AnswerService interface {
	GetById(ctx context.Context, id int64) (models.Answer, error)
	GetList(ctx context.Context, opts AnswerServiceGetListOpts) ([]models.Answer, error)
	GetCount(ctx context.Context, groupId int64) (int64, error)
	Create(ctx context.Context, opts AnswerServiceCreateOpts) (models.Answer, error)
	Update(ctx context.Context, opts AnswerServiceUpdateOpts) (models.Answer, error)
	Delete(ctx context.Context, id int64) error
}

type GroupService interface {
	GetById(ctx context.Context, id int64) (models.Group, error)
	GetList(ctx context.Context, opts GroupServiceGetListOpts) ([]models.Group, error)
	GetCount(ctx context.Context) (int64, error)
	Create(ctx context.Context, opts GroupServiceCreateOpts) (models.Group, error)
	Update(ctx context.Context, opts GroupServiceUpdateOpts) (models.Group, error)
	Delete(ctx context.Context, id int64) error
}

type FileService interface {
	GetById(ctx context.Context, id int64) (models.File, error)
	Create(ctx context.Context, opts FileServiceCreateOpts) (models.File, error)
	Delete(ctx context.Context, id int64) error
}
