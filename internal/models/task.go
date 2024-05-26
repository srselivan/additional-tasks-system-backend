package models

import "time"

type Task struct {
	Id            int64      `json:"id"`
	Title         string     `json:"title"`
	Text          *string    `json:"text"`
	CreatedBy     int64      `json:"createdBy"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
	EffectiveFrom time.Time  `json:"effectiveFrom"`
	EffectiveTill time.Time  `json:"effectiveTill"`
}
