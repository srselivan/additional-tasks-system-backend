package pg

import (
	"backend/internal/models"
	"time"
)

type task struct {
	Id            int64      `db:"id"`
	GroupId       int64      `db:"group_id"`
	Title         string     `db:"title"`
	Text          string     `db:"text"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	EffectiveFrom time.Time  `db:"effective_from"`
	EffectiveTill time.Time  `db:"effective_till"`
	Cost          int64      `db:"cost"`
}

func (t task) toServiceModel() models.Task {
	return models.Task{
		Id:            t.Id,
		GroupId:       t.GroupId,
		Title:         t.Title,
		Text:          t.Text,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
		EffectiveFrom: t.EffectiveFrom,
		EffectiveTill: t.EffectiveTill,
		Cost:          t.Cost,
	}
}
