package answer

import (
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/services/answer/repo"
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type Service struct {
	repo repo.AnswersRepo
	log  *zerolog.Logger
}

func New(
	repo repo.AnswersRepo,
	log *zerolog.Logger,
) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

func (s *Service) GetById(
	ctx context.Context,
	id int64,
) (models.Task, error) {
	task, err := s.repo.GetById(ctx, id)
	if err != nil {
		return models.Task{}, fmt.Errorf("s.repo.GetById: %w", err)
	}
	return task, nil
}

func (s *Service) GetList(
	ctx context.Context,
	opts services.AnswerServiceGetListOpts,
) ([]models.Task, error) {
	tasks, err := s.repo.GetList(ctx, repo.AnswersRepoGetListOpts{
		GroupId: opts.GroupId,
		Limit:   opts.Limit,
		Offset:  opts.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetList: %w", err)
	}
	return tasks, nil
}

func (s *Service) GetCount(
	ctx context.Context,
	groupId int64,
) (int64, error) {
	count, err := s.repo.GetCount(ctx, groupId)
	if err != nil {
		return 0, fmt.Errorf("s.repo.GetCount: %w", err)
	}
	return count, nil
}

func (s *Service) Create(
	ctx context.Context,
	opts services.AnswerServiceCreateOpts,
) (models.Task, error) {
	task, err := s.repo.Create(ctx, repo.AnswersRepoCreateOpts{
		GroupId: opts.GroupId,
		Comment: opts.Comment,
	})
	if err != nil {
		return models.Task{}, fmt.Errorf("s.repo.Create: %w", err)
	}
	return task, nil
}

func (s *Service) Update(
	ctx context.Context,
	opts services.AnswerServiceUpdateOpts,
) (models.Task, error) {
	task, err := s.repo.Update(ctx, repo.AnswersRepoUpdateOpts{
		Id:      opts.Id,
		GroupId: opts.GroupId,
		Comment: opts.Comment,
	})
	if err != nil {
		return models.Task{}, fmt.Errorf("s.repo.Update: %w", err)
	}
	return task, nil
}

func (s *Service) Delete(
	ctx context.Context,
	id int64,
) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("s.repo.Delete: %w", err)
	}
	return nil
}