package services

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"fmt"
	"github.com/samber/lo"
	"time"

	"github.com/rs/zerolog"
)

type TaskService interface {
	GetById(ctx context.Context, id int64) (models.Task, error)
	GetListForCreator(ctx context.Context, opts TaskServiceGetListForCreatorOpts) ([]models.Task, error)
	GetCountForCreator(ctx context.Context, createdBy int64) (int64, error)
	GetListForUser(ctx context.Context, opts TaskServiceGetListForUserOpts) ([]models.Task, error)
	GetCountForUser(ctx context.Context, userId, groupId int64) (int64, error)
	Create(ctx context.Context, opts TaskServiceCreateOpts) (models.Task, error)
	Update(ctx context.Context, opts TaskServiceUpdateOpts) (models.Task, error)
	Delete(ctx context.Context, id int64) error
}

type TaskServiceImpl struct {
	repo             repo.TasksRepo
	filesService     FileService
	taskLinksService TaskLinksService
	log              *zerolog.Logger
}

func NewTaskServiceImpl(
	repo repo.TasksRepo,
	filesService FileService,
	taskLinksService TaskLinksService,
	log *zerolog.Logger,
) *TaskServiceImpl {
	return &TaskServiceImpl{
		repo:             repo,
		log:              log,
		filesService:     filesService,
		taskLinksService: taskLinksService,
	}
}

func (s *TaskServiceImpl) GetById(
	ctx context.Context,
	id int64,
) (models.Task, error) {
	task, err := s.repo.GetById(ctx, id)
	if err != nil {
		return models.Task{}, fmt.Errorf("s.repo.GetById: %w", err)
	}
	return task, nil
}

func (s *TaskServiceImpl) GetListForCreator(
	ctx context.Context,
	opts TaskServiceGetListForCreatorOpts,
) ([]models.Task, error) {
	tasks, err := s.repo.GetListForCreator(ctx, repo.TasksRepoGetListForCreatorOpts{
		CreatedBy: opts.CreatedBy,
		Limit:     opts.Limit,
		Offset:    opts.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetListForCreator: %w", err)
	}
	return tasks, nil
}

func (s *TaskServiceImpl) GetCountForCreator(
	ctx context.Context,
	createdBy int64,
) (int64, error) {
	count, err := s.repo.GetCountForCreator(ctx, createdBy)
	if err != nil {
		return 0, fmt.Errorf("s.repo.GetCountForCreator: %w", err)
	}
	return count, nil
}

func (s *TaskServiceImpl) GetListForUser(
	ctx context.Context,
	opts TaskServiceGetListForUserOpts,
) ([]models.Task, error) {
	tasks, err := s.repo.GetListForUser(ctx, repo.TasksRepoGetListForUserOpts{
		UserId:  opts.UserId,
		GroupId: opts.GroupId,
		Limit:   opts.Limit,
		Offset:  opts.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetListForUser: %w", err)
	}
	return tasks, nil
}

func (s *TaskServiceImpl) GetCountForUser(
	ctx context.Context,
	userId, groupId int64,
) (int64, error) {
	count, err := s.repo.GetCountForUser(ctx, userId, groupId)
	if err != nil {
		return 0, fmt.Errorf("s.repo.GetCountForUser: %w", err)
	}
	return count, nil
}

func (s *TaskServiceImpl) Create(
	ctx context.Context,
	opts TaskServiceCreateOpts,
) (models.Task, error) {
	if opts.EffectiveFrom == nil {
		opts.EffectiveFrom = lo.ToPtr(time.Now())
	}

	task, err := s.repo.Create(ctx, repo.TasksRepoCreateOpts{
		CreatedBy:     opts.CreatedBy,
		Title:         opts.Title,
		Text:          opts.Text,
		EffectiveFrom: *opts.EffectiveFrom,
		EffectiveTill: opts.EffectiveTill,
	})
	if err != nil {
		return models.Task{}, fmt.Errorf("s.repo.Create: %w", err)
	}

	for _, userId := range opts.UserIds {
		if err = s.taskLinksService.Create(ctx, TaskLinksServiceCreateOpts{
			TaskId:  task.Id,
			UserId:  lo.ToPtr(userId),
			GroupId: nil,
		}); err != nil {
			return models.Task{}, fmt.Errorf("s.taskLinksService.Create: %w", err)
		}
	}

	for _, groupId := range opts.GroupIds {
		if err = s.taskLinksService.Create(ctx, TaskLinksServiceCreateOpts{
			TaskId:  task.Id,
			UserId:  nil,
			GroupId: lo.ToPtr(groupId),
		}); err != nil {
			return models.Task{}, fmt.Errorf("s.taskLinksService.Create: %w", err)
		}
	}

	for _, fileId := range opts.FileIds {
		if err = s.filesService.UpdateTaskId(ctx, task.Id, fileId); err != nil {
			return models.Task{}, fmt.Errorf("s.filesService.UpdateTaskId: %d:%w", fileId, err)
		}
	}

	return task, nil
}

func (s *TaskServiceImpl) Update(
	ctx context.Context,
	opts TaskServiceUpdateOpts,
) (models.Task, error) {
	task, err := s.repo.Update(ctx, repo.TasksRepoUpdateOpts{
		Id:            opts.Id,
		GroupId:       opts.GroupId,
		Title:         opts.Title,
		Text:          opts.Text,
		EffectiveFrom: opts.EffectiveFrom,
		EffectiveTill: opts.EffectiveTill,
		Cost:          opts.Cost,
	})
	if err != nil {
		return models.Task{}, fmt.Errorf("s.repo.Update: %w", err)
	}
	return task, nil
}

func (s *TaskServiceImpl) Delete(
	ctx context.Context,
	id int64,
) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("s.repo.Delete: %w", err)
	}
	return nil
}
