package taskshandlers

import (
	"backend/internal/models"
	"time"
)

type getListResponse struct {
	Tasks []models.Task `json:"data"`
	Count int64         `json:"count"`
}

type createRequest struct {
	Title         string    `json:"title"`
	Text          string    `json:"text"`
	EffectiveFrom time.Time `json:"effectiveFrom"`
	EffectiveTill time.Time `json:"effectiveTill"`
	Cost          int64     `json:"cost"`
}

type updateRequest struct {
	Title         string    `json:"title"`
	Text          string    `json:"text"`
	EffectiveFrom time.Time `json:"effectiveFrom"`
	EffectiveTill time.Time `json:"effectiveTill"`
	Cost          int64     `json:"cost"`
}
