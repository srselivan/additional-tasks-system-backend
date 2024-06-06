package repo

import (
	"backend/internal/models"
	"context"
)

type AuthRepo interface {
	SetRefreshToken(ctx context.Context, userId int64, token string) error
	VerifyRefreshToken(ctx context.Context, userId int64, token string) (bool, error)
}

type UsersRepo interface {
	GetById(ctx context.Context, id int64) (models.User, error)
	GetByCredentials(ctx context.Context, opts UsersRepoGetByCredentialsOpts) (models.User, error)
	GetListByRoleId(ctx context.Context, opts UsersRepoGetListByRoleIdOpts) ([]models.User, error)
	GetListByRoleIdCount(ctx context.Context, opts UsersRepoGetListByRoleIdOpts) (int, error)
	Create(ctx context.Context, opts UsersRepoCreateOpts) (models.User, error)
}

type FilesRepo interface {
	GetById(ctx context.Context, id int64) (models.File, error)
	GetByAnswerId(ctx context.Context, answerId int64) ([]models.File, error)
	GetByTaskId(ctx context.Context, taskId int64) ([]models.File, error)
	Create(ctx context.Context, opts FilesRepoCreateOpts) (models.File, error)
	Delete(ctx context.Context, id int64) error
	UpdateAnswerId(ctx context.Context, answerId int64, fileId int64) error
	UpdateTaskId(ctx context.Context, taskId int64, fileId int64) error
}

type AnswersRepo interface {
	GetById(ctx context.Context, id int64) (models.Answer, error)
	GetList(ctx context.Context, opts AnswersRepoGetListOpts) ([]models.Answer, error)
	GetCount(ctx context.Context, groupId int64) (int64, error)
	GetByTaskIdAndUserId(ctx context.Context, userId int64, taskId int64) (models.Answer, error)
	Create(ctx context.Context, opts AnswersRepoCreateOpts) (models.Answer, error)
	Update(ctx context.Context, opts AnswersRepoUpdateOpts) (models.Answer, error)
	Delete(ctx context.Context, id int64) error
}

type MarkRepo interface {
	GetById(ctx context.Context, id int64) (models.Mark, error)
	GetByAnswerId(ctx context.Context, id int64) (models.Mark, error)
	GetListByUserId(ctx context.Context, opts MarksRepoGetListByUserIdOpts) ([]models.Mark, error)
	GetCountByUserId(ctx context.Context, userId int64) (int64, error)
	Create(ctx context.Context, opts MarksRepoCreateOpts) (models.Mark, error)
	Update(ctx context.Context, opts MarksRepoUpdateOpts) (models.Mark, error)
	Delete(ctx context.Context, id int64) error
}

type TasksRepo interface {
	GetById(ctx context.Context, id int64) (models.Task, error)
	GetListForCreator(ctx context.Context, opts TasksRepoGetListForCreatorOpts) ([]models.Task, error)
	GetCountForCreator(ctx context.Context, createdBy int64) (int64, error)
	GetListForUser(ctx context.Context, opts TasksRepoGetListForUserOpts) ([]models.Task, error)
	GetCountForUser(ctx context.Context, userId, groupId int64) (int64, error)
	Create(ctx context.Context, opts TasksRepoCreateOpts) (models.Task, error)
	Update(ctx context.Context, opts TasksRepoUpdateOpts) (models.Task, error)
	Delete(ctx context.Context, id int64) error
}

type TaskLinksRepo interface {
	Create(ctx context.Context, opts TaskLinksRepoCreateOpts) error
}

type GroupsRepo interface {
	GetById(ctx context.Context, id int64) (models.Group, error)
	GetList(ctx context.Context, opts GroupsRepoGetListOpts) ([]models.Group, error)
	GetCount(ctx context.Context) (int64, error)
	Create(ctx context.Context, opts GroupsRepoCreateOpts) (models.Group, error)
	Update(ctx context.Context, opts GroupsRepoUpdateOpts) (models.Group, error)
	Delete(ctx context.Context, id int64) error
}

type StatisticsRepo interface {
	GetStatistics(ctx context.Context, opts GetStatisticsOpts) ([]models.Statistics, error)
}
