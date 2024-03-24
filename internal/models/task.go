package models

import "time"

type Task struct {
	Id            int64      `json:"id"`
	GroupId       int64      `json:"groupId"`
	Title         string     `json:"title"`
	Text          string     `json:"text"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
	EffectiveFrom time.Time  `json:"effectiveFrom"`
	EffectiveTill time.Time  `json:"effectiveTill"`
	Cost          int64      `json:"cost"`
}
