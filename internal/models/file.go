package models

import "time"

type File struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	Filename  string     `json:"filename"`
	Filepath  string     `json:"filepath"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	TaskId    *int64     `json:"taskId"`
	AnswerId  *int64     `json:"answerId"`
}
