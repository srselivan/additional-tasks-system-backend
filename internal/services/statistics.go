package services

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"fmt"
	"github.com/samber/lo"
	"time"
)

var zeroTime = time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)

type StatisticsService interface {
	GetStatistics(ctx context.Context, opts GetStatisticsOpts) ([]models.Statistics, error)
}

type StatisticsServiceImpl struct {
	repo repo.StatisticsRepo
}

func NewStatisticsServiceImpl(repo repo.StatisticsRepo) *StatisticsServiceImpl {
	return &StatisticsServiceImpl{repo: repo}
}

func (s *StatisticsServiceImpl) GetStatistics(
	ctx context.Context,
	opts GetStatisticsOpts,
) ([]models.Statistics, error) {
	if opts.To == nil {
		opts.To = lo.ToPtr(time.Now())
	}
	if opts.From == nil {
		opts.From = lo.ToPtr(zeroTime)
	}

	statistics, err := s.repo.GetStatistics(ctx, repo.GetStatisticsOpts{
		Limit:  opts.Limit,
		Offset: opts.Offset,
		From:   opts.From,
		To:     opts.To,
	})
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetStatistics: %w", err)
	}
	return statistics, nil
}
