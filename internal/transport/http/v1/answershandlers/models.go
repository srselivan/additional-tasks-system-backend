package answershandlers

import (
	"backend/internal/models"
)

type getListResponse struct {
	Answers []models.Answer `json:"data"`
	Count   int64           `json:"count"`
}

type createRequest struct {
	Comment string `json:"comment"`
}

type updateRequest struct {
	Comment string `json:"comment"`
}
