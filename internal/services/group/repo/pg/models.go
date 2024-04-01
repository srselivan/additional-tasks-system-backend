package pg

import (
	"backend/internal/models"
	"time"
)

type group struct {
	Id        int64      `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (g group) toServiceModel() models.Group {
	return models.Group{
		Id:        g.Id,
		Name:      g.Name,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}
