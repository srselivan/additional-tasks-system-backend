package answershandlers

import (
	"backend/internal/models"
)

type getByIdResponse struct {
	Answer        models.Answer `json:"answer"`
	AttachedFiles []models.File `json:"attachedFiles"`
}

type getListResponse struct {
	Data  []extendedAnswer `json:"data"`
	Count int64            `json:"count"`
}

type extendedAnswer struct {
	Answer   models.Answer `json:"answer"`
	UserName string        `json:"userName"`
	Mark     *int64        `json:"mark"`
}

type createRequest struct {
	TaskId  int64   `json:"taskId"`
	Comment *string `json:"comment"`
	Files   []int64 `json:"files"`
}

type updateRequest struct {
	Comment string  `json:"comment"`
	Files   []int64 `json:"files"`
}
