package services

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"time"
)

var (
	ErrUnsuccessfulSignIn   = errors.New("unsuccessful sign in")
	ErrNotFoundRefreshToken = errors.New("such refresh token not exist")
)

type AuthService interface {
	SignIn(ctx context.Context, user models.Credentials) (models.JWTPair, error)
	Refresh(ctx context.Context, refreshToken string) (models.JWTPair, error)
	RefreshTokenExpTime() time.Duration
}

var _ AuthService = (*AuthServiceImpl)(nil)

type AuthServiceImpl struct {
	repo        repo.AuthRepo
	jwtConfig   models.JWTConfig
	userService UserService
	log         *zerolog.Logger
}

func NewAuthServiceImpl(
	repo repo.AuthRepo,
	jwtConfig models.JWTConfig,
	userService UserService,
	log *zerolog.Logger,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		repo:        repo,
		jwtConfig:   jwtConfig,
		userService: userService,
		log:         log,
	}
}

func (s *AuthServiceImpl) SignIn(
	ctx context.Context,
	credentials models.Credentials,
) (models.JWTPair, error) {
	user, err := s.userService.GetByCredentials(ctx, credentials)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return models.JWTPair{}, ErrUnsuccessfulSignIn
		}
		return models.JWTPair{}, fmt.Errorf("s.userService.GetByCredentials: %w", err)
	}

	customClaims := map[string]any{
		"sub":  user.Id,
		"name": user.FullName(),
		"grp":  nil,
		"role": user.RoleId,
	}
	if user.GroupId != nil {
		customClaims["grp"] = *user.GroupId
	}
	jwtPair, err := s.getJWTPair(ctx, customClaims)
	if err != nil {
		return models.JWTPair{}, fmt.Errorf("s.getJWTPair: %w", err)
	}

	if err = s.repo.SetRefreshToken(ctx, user.Id, jwtPair.RefreshToken); err != nil {
		return models.JWTPair{}, fmt.Errorf("s.repo.SetRefreshToken: %w", err)
	}

	return jwtPair, nil
}

func (s *AuthServiceImpl) Refresh(
	ctx context.Context,
	refreshToken string,
) (models.JWTPair, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.JWTRefreshSecretKey), nil
	})
	if err != nil {
		return models.JWTPair{}, fmt.Errorf("jwt.Parse: %w", err)
	}
	claims := token.Claims.(jwt.MapClaims)

	userIdFromClaims, ok := claims["sub"]
	if !ok {
		return models.JWTPair{}, errors.New("not found sub claim in jwt")
	}
	userId, ok := userIdFromClaims.(float64)
	if !ok {
		return models.JWTPair{}, errors.New("sub is not integer")
	}

	ok, err = s.repo.VerifyRefreshToken(ctx, int64(userId), refreshToken)
	if err != nil {
		return models.JWTPair{}, fmt.Errorf("s.repo.VerifyRefreshToken: %w", err)
	}
	if !ok {
		return models.JWTPair{}, ErrNotFoundRefreshToken
	}

	jwtPair, err := s.getJWTPair(ctx, map[string]any{
		"sub":  claims["sub"],
		"name": claims["name"],
		"grp":  claims["grp"],
		"role": claims["role"],
	})
	if err != nil {
		return models.JWTPair{}, fmt.Errorf("s.getJWTPair: %w", err)
	}

	if err = s.repo.SetRefreshToken(ctx, int64(userId), jwtPair.RefreshToken); err != nil {
		return models.JWTPair{}, fmt.Errorf("s.repo.SetRefreshToken: %w", err)
	}

	return jwtPair, nil
}

func (s *AuthServiceImpl) RefreshTokenExpTime() time.Duration {
	return s.jwtConfig.JWTRefreshExpirationTime
}

func (s *AuthServiceImpl) getJWTPair(
	_ context.Context,
	additionalClaims map[string]any,
) (models.JWTPair, error) {
	atClaims := make(jwt.MapClaims, len(additionalClaims)+1)
	atClaims["exp"] = time.Now().Add(s.jwtConfig.JWTAccessExpirationTime).Unix()
	for key, value := range additionalClaims {
		atClaims[key] = value
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(s.jwtConfig.JWTAccessSecretKey))
	if err != nil {
		return models.JWTPair{}, fmt.Errorf("token.SignedString: %w", err)
	}

	rtClaims := make(jwt.MapClaims, len(additionalClaims)+1)
	rtClaims["exp"] = time.Now().Add(s.jwtConfig.JWTRefreshExpirationTime).Unix()
	for key, value := range additionalClaims {
		rtClaims[key] = value
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(s.jwtConfig.JWTRefreshSecretKey))
	if err != nil {
		return models.JWTPair{}, fmt.Errorf("token.SignedString: %w", err)
	}

	return models.JWTPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
