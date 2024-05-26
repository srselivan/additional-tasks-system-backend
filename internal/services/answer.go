package services

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"errors"
	"fmt"
	"github.com/samber/lo"

	"github.com/rs/zerolog"
)

type AnswerService interface {
	GetById(ctx context.Context, id int64) (models.Answer, error)
	GetList(ctx context.Context, opts AnswerServiceGetListOpts) ([]models.Answer, error)
	GetCount(ctx context.Context, taskId int64) (int64, error)
	GetByTaskIdAndUserId(ctx context.Context, userId int64, taskId int64) (models.Answer, error)
	Create(ctx context.Context, opts AnswerServiceCreateOpts) (models.Answer, error)
	Update(ctx context.Context, opts AnswerServiceUpdateOpts) (models.Answer, error)
	Delete(ctx context.Context, id int64) error
}

type AnswerServiceImpl struct {
	repo         repo.AnswersRepo
	filesService FileService
	log          *zerolog.Logger
}

func NewAnswerServiceImpl(
	repo repo.AnswersRepo,
	filesService FileService,
	log *zerolog.Logger,
) *AnswerServiceImpl {
	return &AnswerServiceImpl{
		repo:         repo,
		filesService: filesService,
		log:          log,
	}
}

func (s *AnswerServiceImpl) GetById(
	ctx context.Context,
	id int64,
) (models.Answer, error) {
	answer, err := s.repo.GetById(ctx, id)
	if err != nil {
		return models.Answer{}, fmt.Errorf("s.repo.GetById: %w", err)
	}
	return answer, nil
}

func (s *AnswerServiceImpl) GetList(
	ctx context.Context,
	opts AnswerServiceGetListOpts,
) ([]models.Answer, error) {
	answers, err := s.repo.GetList(ctx, repo.AnswersRepoGetListOpts{
		TaskId: opts.TaskId,
		Limit:  opts.Limit,
		Offset: opts.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetList: %w", err)
	}
	return answers, nil
}

func (s *AnswerServiceImpl) GetByTaskIdAndUserId(
	ctx context.Context,
	userId int64,
	taskId int64,
) (models.Answer, error) {
	answer, err := s.repo.GetByTaskIdAndUserId(ctx, userId, taskId)
	if err != nil {
		return models.Answer{}, fmt.Errorf("s.repo.GetByTaskIdAndUserId: %w", err)
	}
	return answer, nil
}

func (s *AnswerServiceImpl) GetCount(
	ctx context.Context,
	taskId int64,
) (int64, error) {
	count, err := s.repo.GetCount(ctx, taskId)
	if err != nil {
		return 0, fmt.Errorf("s.repo.GetCount: %w", err)
	}
	return count, nil
}

func (s *AnswerServiceImpl) Create(
	ctx context.Context,
	opts AnswerServiceCreateOpts,
) (models.Answer, error) {
	answer, err := s.repo.Create(ctx, repo.AnswersRepoCreateOpts{
		TaskId:  opts.TaskId,
		UserId:  opts.UserId,
		Comment: opts.Comment,
	})
	if err != nil {
		return models.Answer{}, fmt.Errorf("s.repo.Create: %w", err)
	}

	for _, fileId := range opts.Files {
		if err = s.filesService.UpdateAnswerId(ctx, answer.Id, fileId); err != nil {
			return models.Answer{}, fmt.Errorf("s.filesService.UpdateAnswerId: %d:%w", fileId, err)
		}
	}

	return answer, nil
}

func (s *AnswerServiceImpl) Update(
	ctx context.Context,
	opts AnswerServiceUpdateOpts,
) (models.Answer, error) {
	answer, err := s.repo.Update(ctx, repo.AnswersRepoUpdateOpts{
		Id:      opts.Id,
		Comment: opts.Comment,
	})

	files, err := s.filesService.GetByAnswerId(ctx, opts.Id)
	if err != nil {
		return models.Answer{}, fmt.Errorf("s.filesService.GetByAnswerId: %w", err)
	}

	changedFilesMap := lo.SliceToMap(opts.Files, func(item int64) (int64, struct{}) { return item, struct{}{} })

	for _, file := range files {
		if _, ok := changedFilesMap[file.Id]; ok {
			continue
		}

		if err = s.filesService.Delete(ctx, file.Id); err != nil {
			return models.Answer{}, fmt.Errorf("s.filesService.Delete: %w", err)
		}
	}

	for _, fileId := range opts.Files {
		if err = s.filesService.UpdateAnswerId(ctx, opts.Id, fileId); err != nil {
			return models.Answer{}, fmt.Errorf("s.filesService.UpdateAnswerId: %d:%w", fileId, err)
		}
	}

	if err != nil {
		return models.Answer{}, fmt.Errorf("s.repo.Update: %w", err)
	}
	return answer, nil
}

func (s *AnswerServiceImpl) Delete(
	ctx context.Context,
	id int64,
) error {
	files, err := s.filesService.GetByAnswerId(ctx, id)
	if err != nil {
		if !errors.Is(err, repo.ErrNotFound) {
			return fmt.Errorf("s.filesService.GetByAnswerId: %w", err)
		}
	}

	for _, file := range files {
		if err = s.filesService.Delete(ctx, file.Id); err != nil {
			return fmt.Errorf("s.filesService.Delete: %w", err)
		}
	}

	if err = s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("s.repo.Delete: %w", err)
	}
	return nil
}
