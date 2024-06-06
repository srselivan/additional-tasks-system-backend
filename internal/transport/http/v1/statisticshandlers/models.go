package statisticshandlers

import "backend/internal/models"

type getStatisticsResponse struct {
	Data []models.Statistics `json:"data"`
}
