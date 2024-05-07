package models

import "time"

type User struct {
	Id         int64      `json:"id"`
	GroupId    int64      `json:"groupId"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	MiddleName *string    `json:"middleName"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}
