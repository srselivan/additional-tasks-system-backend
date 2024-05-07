package user

import (
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/services/user/repo"
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type Service struct {
	repo repo.UsersRepo
	log  *zerolog.Logger
}

func New(
	repo repo.UsersRepo,
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
) (models.User, error) {
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("s.repo.GetById: %w", err)
	}
	return user, nil
}

func (s *Service) GetByCredentials(
	ctx context.Context,
	opts services.UserServiceGetByCredentialsOpts,
) (models.User, error) {
	user, err := s.repo.GetByCredentials(ctx, repo.UsersRepoGetByCredentialsOpts{})
	if err != nil {
		return models.User{}, fmt.Errorf("s.repo.GetByCredentials: %w", err)
	}
	return user, nil
}

func (s *Service) Create(
	ctx context.Context,
	opts services.UserServiceCreateOpts,
) (models.User, error) {
	user, err := s.repo.Create(ctx, repo.UsersRepoCreateOpts{
		GroupId:    opts.GroupId,
		Email:      opts.Email,
		Password:   opts.Password,
		FirstName:  opts.FirstName,
		LastName:   opts.LastName,
		MiddleName: opts.MiddleName,
	})
	if err != nil {
		return models.User{}, fmt.Errorf("s.repo.Create: %w", err)
	}
	return user, nil
}

func (s *Service) Update(
	ctx context.Context,
	opts services.UserServiceUpdateOpts,
) (models.User, error) {
	//user, err := s.repo.Update(ctx, repo.UsersRepoUpdateOpts{})
	//if err != nil {
	//	return models.User{}, fmt.Errorf("s.repo.Update: %w", err)
	//}
	return models.User{}, nil
}
