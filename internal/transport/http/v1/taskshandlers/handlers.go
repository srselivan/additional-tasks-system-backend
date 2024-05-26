package taskshandlers

import (
	"backend/internal/services"
	"backend/internal/transport/http/auth"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

type handler struct {
	service     services.TaskService
	fileService services.FileService
	log         *zerolog.Logger
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

	attachedFiles, err := h.fileService.GetByTaskId(ctx.UserContext(), task.Id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetById: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(getByIdResponse{
		Task:          task,
		AttachedFiles: attachedFiles,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("jsoniter.Marshal: %v", err))
	}

	if err = ctx.Status(fiber.StatusOK).Send(responseBytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.Send: %v", err))
	}

	return nil
}

func (h *handler) getList(ctx *fiber.Ctx) error {
	claims, err := auth.GetClaimsFromCtx(ctx)
	if err != nil {
		return fmt.Errorf("auth.GetClaimsFromCtx: %w", err)
	}

	limit := ctx.QueryInt("limit")
	if limit == 0 {
		return fiber.NewError(fiber.StatusBadRequest, `Query parameter <limit> missed or equal to zero`)
	}

	offset := ctx.QueryInt("offset", -1)
	if offset == -1 {
		return fiber.NewError(fiber.StatusBadRequest, `Query parameter <offset> missed`)
	}

	creator := ctx.QueryBool("creator")
	if creator {
		tasks, err := h.service.GetListForCreator(ctx.UserContext(), services.TaskServiceGetListForCreatorOpts{
			CreatedBy: claims.UserId,
			Limit:     int64(limit),
			Offset:    int64(offset),
		})
		if err != nil {
			return fmt.Errorf("h.service.GetListForCreator: %w", err)
		}

		count, err := h.service.GetCountForCreator(ctx.UserContext(), claims.UserId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetCountForCreator: %v", err))
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

	tasks, err := h.service.GetListForUser(ctx.UserContext(), services.TaskServiceGetListForUserOpts{
		UserId:  claims.UserId,
		GroupId: lo.FromPtr(claims.GroupId),
		Limit:   int64(limit),
		Offset:  int64(offset),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetListForUser: %v", err))
	}

	count, err := h.service.GetCountForUser(ctx.UserContext(), claims.UserId, lo.FromPtr(claims.GroupId))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetCountForUser: %v", err))
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
	claims, err := auth.GetClaimsFromCtx(ctx)
	if err != nil {
		return fmt.Errorf("auth.GetClaimsFromCtx: %w", err)
	}

	var req createRequest
	if err = jsoniter.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("jsoniter.Unmarshal: %v", err))
	}

	task, err := h.service.Create(ctx.UserContext(), services.TaskServiceCreateOpts{
		GroupIds:      req.GroupIds,
		UserIds:       req.UserIds,
		Title:         req.Title,
		Text:          req.Text,
		EffectiveFrom: req.EffectiveFrom,
		EffectiveTill: req.EffectiveTill,
		FileIds:       req.FileIds,
		CreatedBy:     claims.UserId,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.Create: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(task)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("jsoniter.Marshal: %v", err))
	}

	if err = ctx.Status(fiber.StatusCreated).Send(responseBytes); err != nil {
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
