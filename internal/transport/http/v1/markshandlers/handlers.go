package markshandlers

import (
	"backend/internal/services"
	"backend/internal/transport/http/auth"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

type handler struct {
	service services.MarkService
	log     *zerolog.Logger
}

func (h *handler) getById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, `Path parameter <id> empty or not a number`)
	}

	mark, err := h.service.GetById(ctx.UserContext(), int64(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetById: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(mark)
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

	offset := ctx.QueryInt("offset", -1)
	if offset == -1 {
		return fiber.NewError(fiber.StatusBadRequest, `Query parameter <offset> missed`)
	}

	taskId := ctx.QueryInt("taskId", -1)
	if taskId != -1 {
		return nil
	}

	claims, err := auth.GetClaimsFromCtx(ctx)
	if err != nil {
		return fmt.Errorf("auth.GetClaimsFromCtx: %v", err)
	}

	marks, err := h.service.GetListByUserId(ctx.UserContext(), services.MarkServiceGetListByUserIdOpts{
		UserId: claims.UserId,
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetListByUserId: %v", err))
	}

	count, err := h.service.GetCountByUserId(ctx.UserContext(), claims.UserId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetCountByUserId: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(getListResponse{
		Marks: marks,
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
	var req createMarkRequest
	if err := jsoniter.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("jsoniter.Unmarshal: %v", err))
	}

	mark, err := h.service.Create(ctx.UserContext(), services.MarkServiceCreateOpts{
		AnswerId: req.AnswerId,
		Mark:     req.Mark,
		Comment:  req.Comment,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.Create: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(mark)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("jsoniter.Marshal: %v", err))
	}

	if err = ctx.Status(fiber.StatusCreated).Send(responseBytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.Send: %v", err))
	}

	return nil
}

func (h *handler) update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, `Path parameter <id> empty or not a number`)
	}

	var req updateMarkRequest
	if err = jsoniter.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("jsoniter.Unmarshal: %v", err))
	}

	mark, err := h.service.Update(ctx.UserContext(), services.MarkServiceUpdateOpts{
		Id:      int64(id),
		Mark:    req.Mark,
		Comment: req.Comment,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.Update: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(mark)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("jsoniter.Marshal: %v", err))
	}

	if err = ctx.Status(fiber.StatusCreated).Send(responseBytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.Send: %v", err))
	}

	return nil
}
