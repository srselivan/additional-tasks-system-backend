package services

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	Create(ctx context.Context, opts UserServiceCreateOpts) (models.User, error)
	GetById(ctx context.Context, id int64) (models.User, error)
	GetByCredentials(ctx context.Context, credentials models.Credentials) (models.User, error)
	GetListByRoleId(ctx context.Context, opts UserServiceGetListByRoleIdOpts) ([]models.User, error)
	GetListByRoleIdCount(ctx context.Context, opts UserServiceGetListByRoleIdOpts) (int, error)
}

type UserServiceImpl struct {
	repo repo.UsersRepo
	log  *zerolog.Logger
}

func NewUserServiceImpl(
	repo repo.UsersRepo,
	log *zerolog.Logger,
) *UserServiceImpl {
	return &UserServiceImpl{
		repo: repo,
		log:  log,
	}
}

func (s *UserServiceImpl) GetById(
	ctx context.Context,
	id int64,
) (models.User, error) {
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("s.repo.GetById: %w", err)
	}
	return user, nil
}

func (s *UserServiceImpl) GetByCredentials(
	ctx context.Context,
	credentials models.Credentials,
) (models.User, error) {
	user, err := s.repo.GetByCredentials(ctx, repo.UsersRepoGetByCredentialsOpts{
		Email:    credentials.Email,
		Password: credentials.Password,
	})
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("s.repo.GetByCredentials: %w", err)
	}
	return user, nil
}

func (s *UserServiceImpl) GetListByRoleId(
	ctx context.Context,
	opts UserServiceGetListByRoleIdOpts,
) ([]models.User, error) {
	users, err := s.repo.GetListByRoleId(ctx, repo.UsersRepoGetListByRoleIdOpts{
		RoleId: opts.RoleId,
		Limit:  opts.Limit,
		Offset: opts.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetListByRoleId: %w", err)
	}
	return users, nil
}

func (s *UserServiceImpl) GetListByRoleIdCount(
	ctx context.Context,
	opts UserServiceGetListByRoleIdOpts,
) (int, error) {
	count, err := s.repo.GetListByRoleIdCount(ctx, repo.UsersRepoGetListByRoleIdOpts{
		RoleId: opts.RoleId,
		Limit:  opts.Limit,
		Offset: opts.Offset,
	})
	if err != nil {
		return 0, fmt.Errorf("s.repo.GetListByRoleIdCount: %w", err)
	}
	return count, nil
}

type UserServiceCreateOpts struct {
	GroupId    *int64
	RoleId     int64
	Email      string
	Password   string
	FirstName  string
	LastName   string
	MiddleName *string
}

func (s *UserServiceImpl) Create(
	ctx context.Context,
	opts UserServiceCreateOpts,
) (models.User, error) {
	user, err := s.repo.Create(ctx, repo.UsersRepoCreateOpts{
		GroupId:    opts.GroupId,
		RoleId:     opts.RoleId,
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

type UserServiceUpdateOpts struct {
	GroupId    int64
	Email      string
	Password   string
	FirstName  string
	LastName   string
	MiddleName *string
}

func (s *UserServiceImpl) Update(
	ctx context.Context,
	opts UserServiceUpdateOpts,
) (models.User, error) {
	//user, err := s.repo.Update(ctx, repo.UsersRepoUpdateOpts{})
	//if err != nil {
	//	return models.User{}, fmt.Errorf("s.repo.Update: %w", err)
	//}
	return models.User{}, nil
}
