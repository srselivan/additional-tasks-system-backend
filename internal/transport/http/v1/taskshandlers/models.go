package taskshandlers

import (
	"backend/internal/models"
	"time"
)

type getByIdResponse struct {
	Task          models.Task   `json:"task"`
	AttachedFiles []models.File `json:"attachedFiles"`
}

type getListResponse struct {
	Tasks []models.Task `json:"data"`
	Count int64         `json:"count"`
}

type createRequest struct {
	GroupIds      []int64    `json:"groupIds"`
	UserIds       []int64    `json:"userIds"`
	Title         string     `json:"title"`
	Text          *string    `json:"text,omitempty"`
	EffectiveFrom *time.Time `json:"effectiveFrom"`
	EffectiveTill time.Time  `json:"effectiveTill"`
	FileIds       []int64    `json:"fileIds"`
}

type updateRequest struct {
	Title         string    `json:"title"`
	Text          string    `json:"text"`
	EffectiveFrom time.Time `json:"effectiveFrom"`
	EffectiveTill time.Time `json:"effectiveTill"`
	Cost          int64     `json:"cost"`
}
