package services

import (
	"backend/internal/repo"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type TaskLinksService interface {
	Create(ctx context.Context, opts TaskLinksServiceCreateOpts) error
}

type TaskLinksServiceImpl struct {
	repo repo.TaskLinksRepo
	log  *zerolog.Logger
}

func NewTaskLinksServiceImpl(
	repo repo.TaskLinksRepo,
	log *zerolog.Logger,
) *TaskLinksServiceImpl {
	return &TaskLinksServiceImpl{
		repo: repo,
		log:  log,
	}
}

func (s *TaskLinksServiceImpl) Create(
	ctx context.Context,
	opts TaskLinksServiceCreateOpts,
) error {
	if err := s.repo.Create(ctx, repo.TaskLinksRepoCreateOpts{
		TaskId:  opts.TaskId,
		UserId:  opts.UserId,
		GroupId: opts.GroupId,
	}); err != nil {
		return fmt.Errorf("s.repo.Create: %w", err)
	}
	return nil
}
