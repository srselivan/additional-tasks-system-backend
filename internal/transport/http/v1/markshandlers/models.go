package markshandlers

import (
	"backend/internal/models"
)

type getListResponse struct {
	Marks []models.Mark `json:"data"`
	Count int64         `json:"count"`
}

type createMarkRequest struct {
	AnswerId int64   `json:"answerId"`
	Mark     int64   `json:"mark"`
	Comment  *string `json:"comment"`
}

type updateMarkRequest struct {
	Mark    int64   `json:"mark"`
	Comment *string `json:"comment"`
}
