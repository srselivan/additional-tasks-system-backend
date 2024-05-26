package usershandlers

import (
	"backend/internal/models"
)

type getListResponse struct {
	Users []models.User `json:"data"`
	Count int           `json:"count"`
}
