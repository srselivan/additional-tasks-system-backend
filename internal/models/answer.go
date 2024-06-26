package models

import "time"

type Answer struct {
	Id        int64      `json:"id"`
	UserId    int64      `json:"userId"`
	TaskId    int64      `json:"taskId"`
	Comment   *string    `json:"comment,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
