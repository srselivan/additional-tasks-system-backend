package groupshandlers

import (
	"backend/internal/models"
)

type getListResponse struct {
	Groups []models.Group `json:"data"`
	Count  int64          `json:"count"`
}

type createRequest struct {
	Name string `json:"name"`
}

type updateRequest struct {
	Name string `json:"name"`
}
