package usershandlers

import (
	"backend/internal/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

type handler struct {
	service services.UserService
	log     *zerolog.Logger
}

func (h *handler) getById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, `Path parameter <id> empty or not a number`)
	}

	user, err := h.service.GetById(ctx.UserContext(), int64(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetById: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(user)
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

	roleId := ctx.QueryInt("roleId", -1)
	if roleId == -1 {
		return fiber.NewError(fiber.StatusBadRequest, `Query parameter <roleId> missed`)
	}

	users, err := h.service.GetListByRoleId(ctx.UserContext(), services.UserServiceGetListByRoleIdOpts{
		RoleId: int64(roleId),
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetListByRoleId: %v", err))
	}

	count, err := h.service.GetListByRoleIdCount(ctx.UserContext(), services.UserServiceGetListByRoleIdOpts{
		RoleId: int64(roleId),
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetListByRoleIdCount: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(getListResponse{
		Users: users,
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
