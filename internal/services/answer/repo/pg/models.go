package pg

import (
	"backend/internal/models"
	"time"
)

type answer struct {
	Id        int64      `db:"id"`
	GroupId   int64      `db:"group_id"`
	Comment   *string    `db:"comment"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (a answer) toServiceModel() models.Answer {
	return models.Answer{
		Id:        a.Id,
		GroupId:   a.GroupId,
		Comment:   a.Comment,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
