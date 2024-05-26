package services

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type MarkService interface {
	GetById(ctx context.Context, id int64) (models.Mark, error)
	GetByAnswerId(ctx context.Context, id int64) (models.Mark, error)
	GetListByUserId(ctx context.Context, opts MarkServiceGetListByUserIdOpts) ([]models.Mark, error)
	GetCountByUserId(ctx context.Context, userId int64) (int64, error)
	Create(ctx context.Context, opts MarkServiceCreateOpts) (models.Mark, error)
	Update(ctx context.Context, opts MarkServiceUpdateOpts) (models.Mark, error)
	Delete(ctx context.Context, id int64) error
}

type MarkServiceImpl struct {
	repo repo.MarkRepo
	log  *zerolog.Logger
}

func NewMarkServiceImpl(
	repo repo.MarkRepo,
	log *zerolog.Logger,
) *MarkServiceImpl {
	return &MarkServiceImpl{repo: repo, log: log}
}

func (s *MarkServiceImpl) GetById(ctx context.Context, id int64) (models.Mark, error) {
	mark, err := s.repo.GetById(ctx, id)
	if err != nil {
		return models.Mark{}, fmt.Errorf("s.repo.GetById: %w", err)
	}
	return mark, nil
}

func (s *MarkServiceImpl) GetByAnswerId(ctx context.Context, id int64) (models.Mark, error) {
	mark, err := s.repo.GetByAnswerId(ctx, id)
	if err != nil {
		return models.Mark{}, fmt.Errorf("s.repo.GetByAnswerId: %w", err)
	}
	return mark, nil
}

func (s *MarkServiceImpl) GetListByUserId(
	ctx context.Context,
	opts MarkServiceGetListByUserIdOpts,
) ([]models.Mark, error) {
	marks, err := s.repo.GetListByUserId(ctx, repo.MarksRepoGetListByUserIdOpts{
		Limit:  opts.Limit,
		Offset: opts.Offset,
		UserId: opts.UserId,
	})
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetListByUserId: %w", err)
	}
	return marks, nil
}

func (s *MarkServiceImpl) GetCountByUserId(
	ctx context.Context,
	userId int64,
) (int64, error) {
	count, err := s.repo.GetCountByUserId(ctx, userId)
	if err != nil {
		return 0, fmt.Errorf("s.repo.GetCountByUserId: %w", err)
	}
	return count, nil
}

func (s *MarkServiceImpl) Create(
	ctx context.Context,
	opts MarkServiceCreateOpts,
) (models.Mark, error) {
	mark, err := s.repo.Create(ctx, repo.MarksRepoCreateOpts{
		AnswerId: opts.AnswerId,
		Mark:     opts.Mark,
		Comment:  opts.Comment,
	})
	if err != nil {
		return models.Mark{}, fmt.Errorf("s.repo.Create: %w", err)
	}
	return mark, nil
}

func (s *MarkServiceImpl) Update(
	ctx context.Context,
	opts MarkServiceUpdateOpts,
) (models.Mark, error) {
	mark, err := s.repo.Update(ctx, repo.MarksRepoUpdateOpts{
		Id:      opts.Id,
		Mark:    opts.Mark,
		Comment: opts.Comment,
	})
	if err != nil {
		return models.Mark{}, fmt.Errorf("s.repo.Update: %w", err)
	}
	return mark, nil
}

func (s *MarkServiceImpl) Delete(
	ctx context.Context,
	id int64,
) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("s.repo.Delete: %w", err)
	}
	return nil
}
