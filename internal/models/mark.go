package models

import "time"

type Mark struct {
	Id        int64      `json:"id"`
	AnswerId  int64      `json:"answerId"`
	TaskId    int64      `json:"taskId"`
	Mark      int64      `json:"mark"`
	Comment   *string    `json:"comment,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
