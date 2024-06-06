package statisticshandlers

import (
	"backend/internal/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	"time"
)

type handler struct {
	service services.StatisticsService
	log     *zerolog.Logger
}

func (h *handler) get(ctx *fiber.Ctx) error {
	limit := ctx.QueryInt("limit")
	if limit == 0 {
		return fiber.NewError(fiber.StatusBadRequest, `Query parameter <limit> missed or equal to zero`)
	}

	offset := ctx.QueryInt("offset", -1)
	if offset == -1 {
		return fiber.NewError(fiber.StatusBadRequest, `Query parameter <offset> missed`)
	}

	opts := services.GetStatisticsOpts{
		Limit:  int64(limit),
		Offset: int64(offset),
	}

	queries := ctx.Queries()

	from, ok := queries["from"]
	if ok {
		fromTime, err := time.Parse("2006-01-02T15:04:05Z", from)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, `Query parameter <from> incorrect format`)
		}
		opts.From = &fromTime
	}

	to, ok := queries["to"]
	if ok {
		toTime, err := time.Parse("2006-01-02T15:04:05Z", to)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, `Query parameter <to> incorrect format`)
		}
		opts.To = &toTime
	}

	statistics, err := h.service.GetStatistics(ctx.UserContext(), opts)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetStatistics: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(getStatisticsResponse{
		Data: statistics,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("jsoniter.Marshal: %v", err))
	}

	if err = ctx.Status(fiber.StatusOK).Send(responseBytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.Send: %v", err))
	}

	return nil
}
