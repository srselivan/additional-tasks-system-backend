package services

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type FileService interface {
	GetById(ctx context.Context, id int64) (models.File, error)
	GetByAnswerId(ctx context.Context, answerId int64) ([]models.File, error)
	GetByTaskId(ctx context.Context, taskId int64) ([]models.File, error)
	Create(ctx context.Context, opts FileServiceCreateOpts) (models.File, error)
	UpdateAnswerId(ctx context.Context, answerId int64, fileId int64) error
	UpdateTaskId(ctx context.Context, taskId int64, fileId int64) error
	Delete(ctx context.Context, id int64) error
}

type FileServiceImpl struct {
	repo repo.FilesRepo
	log  *zerolog.Logger
}

func NewFileServiceImpl(
	repo repo.FilesRepo,
	log *zerolog.Logger,
) *FileServiceImpl {
	return &FileServiceImpl{
		repo: repo,
		log:  log,
	}
}

func (s *FileServiceImpl) GetById(
	ctx context.Context,
	id int64,
) (models.File, error) {
	file, err := s.repo.GetById(ctx, id)
	if err != nil {
		return models.File{}, fmt.Errorf("s.repo.GetById: %w", err)
	}
	return file, nil
}

func (s *FileServiceImpl) GetByAnswerId(
	ctx context.Context,
	answerId int64,
) ([]models.File, error) {
	files, err := s.repo.GetByAnswerId(ctx, answerId)
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetByAnswerId: %w", err)
	}
	return files, nil
}

func (s *FileServiceImpl) GetByTaskId(
	ctx context.Context,
	taskId int64,
) ([]models.File, error) {
	files, err := s.repo.GetByTaskId(ctx, taskId)
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetByTaskId: %w", err)
	}
	return files, nil
}

func (s *FileServiceImpl) Create(
	ctx context.Context,
	opts FileServiceCreateOpts,
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

func (s *FileServiceImpl) UpdateAnswerId(
	ctx context.Context,
	answerId int64,
	fileId int64,
) error {
	if err := s.repo.UpdateAnswerId(ctx, answerId, fileId); err != nil {
		return fmt.Errorf("s.repo.UpdateAnswerId: %w", err)
	}
	return nil
}

func (s *FileServiceImpl) UpdateTaskId(
	ctx context.Context,
	taskId int64,
	fileId int64,
) error {
	if err := s.repo.UpdateTaskId(ctx, taskId, fileId); err != nil {
		return fmt.Errorf("s.repo.UpdateTaskId: %w", err)
	}
	return nil
}

func (s *FileServiceImpl) Delete(
	ctx context.Context,
	id int64,
) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("s.repo.Delete: %w", err)
	}
	return nil
}
