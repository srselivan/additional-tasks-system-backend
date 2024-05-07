package pg

import (
	"backend/internal/models"
	"time"
)

type user struct {
	Id         int64      `db:"id"`
	GroupId    int64      `db:"group_id"`
	Email      string     `db:"email"`
	Password   string     `db:"password"`
	FirstName  string     `db:"first_name"`
	LastName   string     `db:"last_name"`
	MiddleName *string    `db:"middle_name"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

func (u user) toServiceModel() models.User {
	return models.User{
		Id:         u.Id,
		GroupId:    u.GroupId,
		Email:      u.Email,
		Password:   u.Password,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		MiddleName: u.MiddleName,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}
