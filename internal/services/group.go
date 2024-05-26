package services

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type GroupService interface {
	GetById(ctx context.Context, id int64) (models.Group, error)
	GetList(ctx context.Context, opts GroupServiceGetListOpts) ([]models.Group, error)
	GetCount(ctx context.Context) (int64, error)
	Create(ctx context.Context, opts GroupServiceCreateOpts) (models.Group, error)
	Update(ctx context.Context, opts GroupServiceUpdateOpts) (models.Group, error)
	Delete(ctx context.Context, id int64) error
}

type GroupServiceImpl struct {
	repo repo.GroupsRepo
	log  *zerolog.Logger
}

func NewGroupServiceImpl(
	repo repo.GroupsRepo,
	log *zerolog.Logger,
) *GroupServiceImpl {
	return &GroupServiceImpl{
		repo: repo,
		log:  log,
	}
}

func (s *GroupServiceImpl) GetById(
	ctx context.Context,
	id int64,
) (models.Group, error) {
	group, err := s.repo.GetById(ctx, id)
	if err != nil {
		return models.Group{}, fmt.Errorf("s.repo.GetById: %w", err)
	}
	return group, nil
}

func (s *GroupServiceImpl) GetList(
	ctx context.Context,
	opts GroupServiceGetListOpts,
) ([]models.Group, error) {
	groups, err := s.repo.GetList(ctx, repo.GroupsRepoGetListOpts{
		Limit:  opts.Limit,
		Offset: opts.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetList: %w", err)
	}
	return groups, nil
}

func (s *GroupServiceImpl) GetCount(
	ctx context.Context,
) (int64, error) {
	count, err := s.repo.GetCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("s.repo.GetCount: %w", err)
	}
	return count, nil
}

func (s *GroupServiceImpl) Create(
	ctx context.Context,
	opts GroupServiceCreateOpts,
) (models.Group, error) {
	group, err := s.repo.Create(ctx, repo.GroupsRepoCreateOpts{
		Name: opts.Name,
	})
	if err != nil {
		return models.Group{}, fmt.Errorf("s.repo.Create: %w", err)
	}
	return group, nil
}

func (s *GroupServiceImpl) Update(
	ctx context.Context,
	opts GroupServiceUpdateOpts,
) (models.Group, error) {
	group, err := s.repo.Update(ctx, repo.GroupsRepoUpdateOpts{
		Id:   opts.Id,
		Name: opts.Name,
	})
	if err != nil {
		return models.Group{}, fmt.Errorf("s.repo.Update: %w", err)
	}
	return group, nil
}

func (s *GroupServiceImpl) Delete(
	ctx context.Context,
	id int64,
) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("s.repo.Delete: %w", err)
	}
	return nil
}
