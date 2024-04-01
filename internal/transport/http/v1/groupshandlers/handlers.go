package groupshandlers

import (
	"backend/internal/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

type handler struct {
	service services.GroupService
	log     *zerolog.Logger
}

func (h *handler) getById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, `Path parameter <id> empty or not a number`)
	}

	group, err := h.service.GetById(ctx.UserContext(), int64(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetById: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(group)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("jsoniter.Marshal: %v", err))
	}

	if err = ctx.Status(fiber.StatusOK).Send(responseBytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.Send: %v", err))
	}

	return nil
}

func (h *handler) getList(ctx *fiber.Ctx) error {
	limit := ctx.QueryInt("limit")
	if limit == 0 {
		return fiber.NewError(fiber.StatusBadRequest, `Query parameter <limit> missed or equal to zero`)
	}

	offset := ctx.QueryInt("offset")
	if offset == 0 {
		return fiber.NewError(fiber.StatusBadRequest, `Query parameter <offset> missed or equal to zero`)
	}

	groups, err := h.service.GetList(ctx.UserContext(), services.GroupServiceGetListOpts{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetList: %v", err))
	}

	count, err := h.service.GetCount(ctx.UserContext())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetCount: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(getListResponse{
		Groups: groups,
		Count:  count,
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
	var req createRequest
	if err := jsoniter.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("jsoniter.Unmarshal: %v", err))
	}

	_, err := h.service.Create(ctx.UserContext(), services.GroupServiceCreateOpts{
		Name: req.Name,
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
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, `Path parameter <id> empty or not a number`)
	}

	var req updateRequest
	if err = jsoniter.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("jsoniter.Unmarshal: %v", err))
	}

	_, err = h.service.Update(ctx.UserContext(), services.GroupServiceUpdateOpts{
		Id:   int64(id),
		Name: req.Name,
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
