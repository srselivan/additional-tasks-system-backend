package models

import "time"

type Answer struct {
	Id        int64      `json:"id"`
	GroupId   int64      `json:"groupId"`
	Comment   *string    `json:"comment"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
