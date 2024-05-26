package models

import (
	"time"
)

const (
	UserRoleAdministrator = 1
	UserRoleTeacher       = 2
	UserRoleStudent       = 3
)

type User struct {
	Id         int64      `json:"id"`
	GroupId    *int64     `json:"groupId,omitempty"`
	RoleId     int64      `json:"roleId"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	MiddleName *string    `json:"middleName"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}

func (u User) FullName() string {
	if u.MiddleName == nil {
		return u.LastName + " " + u.FirstName
	}
	return u.LastName + " " + u.FirstName + " " + *u.MiddleName
}
