package pg

import (
	"backend/internal/models"
	"backend/internal/repo"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type statistics struct {
	UserName  string `db:"user_name"`
	GroupName string `db:"group_name"`
	Score     int64  `db:"score"`
}

type StatisticsRepo struct {
	db *sqlx.DB
}

func NewStatisticsRepo(db *sqlx.DB) *StatisticsRepo {
	return &StatisticsRepo{db: db}
}

const statisticsRepoGetStatisticsQuery = `
select 
	u.last_name || ' ' || u.first_name || ' ' || coalesce(u.middle_name, '') user_name,
	g.name group_name,
	sum(m.mark) score
from public.mark m
left join public.answer a on a.id = m.answer_id
left join public."user" u on a.user_id = u.id
left join public."group" g on g.id = u.group_id 
where coalesce(m.updated_at, m.created_at) > $1 and coalesce(m.updated_at, m.created_at) < $2
group by u.last_name, u.middle_name, u.first_name, g.name
order by score desc 
limit $3
offset $4
`

func (r *StatisticsRepo) GetStatistics(
	ctx context.Context,
	opts repo.GetStatisticsOpts,
) ([]models.Statistics, error) {
	var s []statistics
	if err := r.db.SelectContext(
		ctx,
		&s,
		statisticsRepoGetStatisticsQuery,
		opts.From,
		opts.To,
		opts.Limit,
		opts.Offset,
	); err != nil {
		return nil, fmt.Errorf("r.db.SelectContext: %w", err)
	}
	return lo.Map(
		s,
		func(item statistics, _ int) models.Statistics {
			return models.Statistics{
				UserName:  item.UserName,
				GroupName: item.GroupName,
				Score:     item.Score,
			}
		},
	), nil
}
