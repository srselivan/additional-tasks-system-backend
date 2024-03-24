package taskshandlers

import (
	"backend/internal/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

type handler struct {
	service services.TaskService
	log     *zerolog.Logger
}

func (h *handler) getById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, `Path parameter <id> empty or not a number`)
	}

	task, err := h.service.GetById(ctx.UserContext(), int64(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetById: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(task)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("jsoniter.Marshal: %v", err))
	}

	if err = ctx.Status(fiber.StatusOK).Send(responseBytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.Send: %v", err))
	}

	return nil
}

func (h *handler) getList(ctx *fiber.Ctx) error {
	groupId := 0

	limit := ctx.QueryInt("limit")
	if limit == 0 {
		return fiber.NewError(fiber.StatusBadRequest, `Query parameter <limit> missed or equal to zero`)
	}

	offset := ctx.QueryInt("offset")
	if offset == 0 {
		return fiber.NewError(fiber.StatusBadRequest, `Query parameter <offset> missed or equal to zero`)
	}

	tasks, err := h.service.GetList(ctx.UserContext(), services.TaskServiceGetListOpts{
		GroupId: int64(groupId),
		Limit:   int64(limit),
		Offset:  int64(offset),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetList: %v", err))
	}

	count, err := h.service.GetCount(ctx.UserContext(), int64(groupId))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetCount: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(getListResponse{
		Tasks: tasks,
		Count: count,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("jsoniter.Marshal: %v", err))
	}

	if err = ctx.Status(fiber.StatusOK).Send(responseBytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.Send: %v", err))
	}

	return nil
}

func (h *handler) create(ctx *fiber.Ctx) error {
	groupId := 0

	var req createRequest
	if err := jsoniter.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("jsoniter.Unmarshal: %v", err))
	}

	_, err := h.service.Create(ctx.UserContext(), services.TaskServiceCreateOpts{
		GroupId:       int64(groupId),
		Title:         req.Title,
		Text:          req.Text,
		EffectiveFrom: req.EffectiveFrom,
		EffectiveTill: req.EffectiveTill,
		Cost:          req.Cost,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.Create: %v", err))
	}

	if err = ctx.Status(fiber.StatusCreated).Send(nil); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.Send: %v", err))
	}

	return nil
}

func (h *handler) update(ctx *fiber.Ctx) error {
	groupId := 0

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, `Path parameter <id> empty or not a number`)
	}

	var req updateRequest
	if err = jsoniter.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("jsoniter.Unmarshal: %v", err))
	}

	_, err = h.service.Update(ctx.UserContext(), services.TaskServiceUpdateOpts{
		Id:            int64(id),
		GroupId:       int64(groupId),
		Title:         req.Title,
		Text:          req.Text,
		EffectiveFrom: req.EffectiveFrom,
		EffectiveTill: req.EffectiveTill,
		Cost:          req.Cost,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.Create: %v", err))
	}

	if err = ctx.Status(fiber.StatusAccepted).Send(nil); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.Send: %v", err))
	}

	return nil
}

func (h *handler) delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, `Path parameter <id> empty or not a number`)
	}

	if err = h.service.Delete(ctx.UserContext(), int64(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.Delete: %v", err))
	}

	if err = ctx.Status(fiber.StatusAccepted).Send(nil); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.Send: %v", err))
	}

	return nil
}
