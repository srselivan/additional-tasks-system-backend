package file

import (
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/services/file/repo"
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type Service struct {
	repo repo.FilesRepo
	log  *zerolog.Logger
}

func New(
	repo repo.FilesRepo,
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
) (models.File, error) {
	file, err := s.repo.GetById(ctx, id)
	if err != nil {
		return models.File{}, fmt.Errorf("s.repo.GetById: %w", err)
	}
	return file, nil
}

func (s *Service) Create(
	ctx context.Context,
	opts services.FileServiceCreateOpts,
) (models.File, error) {
	file, err := s.repo.Create(ctx, repo.FilesRepoCreateOpts{
		Name:     opts.Name,
		Filename: opts.Filename,
		Filepath: opts.Filepath,
		TaskId:   opts.TaskId,
		AnswerId: opts.AnswerId,
	})
	if err != nil {
		return models.File{}, fmt.Errorf("s.repo.Create: %w", err)
	}
	return file, nil
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
