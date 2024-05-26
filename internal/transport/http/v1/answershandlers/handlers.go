package answershandlers

import (
	"backend/internal/models"
	"backend/internal/repo"
	"backend/internal/services"
	"backend/internal/transport/http/auth"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

type handler struct {
	service     services.AnswerService
	fileService services.FileService
	userService services.UserService
	markService services.MarkService
	log         *zerolog.Logger
}

func (h *handler) getById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, `Path parameter <id> empty or not a number`)
	}

	answer, err := h.service.GetById(ctx.UserContext(), int64(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetById: %v", err))
	}

	files, err := h.fileService.GetByAnswerId(ctx.UserContext(), answer.Id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetById: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(getByIdResponse{
		Answer:        answer,
		AttachedFiles: files,
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

	taskId := ctx.QueryInt("taskId", -1)
	if taskId == -1 {
		return fiber.NewError(fiber.StatusBadRequest, `Query parameter <taskId> missed`)
	}

	if claims.Role == models.UserRoleStudent {
		answer, err := h.service.GetByTaskIdAndUserId(ctx.UserContext(), claims.UserId, int64(taskId))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetByTaskId: %v", err))
		}

		files, err := h.fileService.GetByAnswerId(ctx.UserContext(), answer.Id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetById: %v", err))
		}

		responseBytes, err := jsoniter.Marshal(getByIdResponse{
			Answer:        answer,
			AttachedFiles: files,
		})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("jsoniter.Marshal: %v", err))
		}

		if err = ctx.Status(fiber.StatusOK).Send(responseBytes); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.Send: %v", err))
		}

		return nil
	}

	limit := ctx.QueryInt("limit")
	if limit == 0 {
		return fiber.NewError(fiber.StatusBadRequest, `Query parameter <limit> missed or equal to zero`)
	}

	offset := ctx.QueryInt("offset", -1)
	if offset == -1 {
		return fiber.NewError(fiber.StatusBadRequest, `Query parameter <offset> missed`)
	}

	answers, err := h.service.GetList(ctx.UserContext(), services.AnswerServiceGetListOpts{
		TaskId: int64(taskId),
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetList: %v", err))
	}

	count, err := h.service.GetCount(ctx.UserContext(), int64(taskId))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetCount: %v", err))
	}

	data := make([]extendedAnswer, 0, len(answers))
	for _, answer := range answers {
		user, err := h.userService.GetById(ctx.UserContext(), answer.UserId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.userService.GetById: %v", err))
		}

		mark, err := h.markService.GetByAnswerId(ctx.UserContext(), answer.Id)
		if err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				data = append(data, extendedAnswer{
					Answer:   answer,
					UserName: user.FullName(),
					Mark:     nil,
				})
				continue
			}
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.markService.GetById: %v", err))
		}

		data = append(data, extendedAnswer{
			Answer:   answer,
			UserName: user.FullName(),
			Mark:     &mark.Mark,
		})
	}

	responseBytes, err := jsoniter.Marshal(getListResponse{
		Data:  data,
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

	_, err = h.service.Create(ctx.UserContext(), services.AnswerServiceCreateOpts{
		TaskId:  req.TaskId,
		UserId:  claims.UserId,
		Comment: req.Comment,
		Files:   req.Files,
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

	_, err = h.service.Update(ctx.UserContext(), services.AnswerServiceUpdateOpts{
		Id:      int64(id),
		Comment: req.Comment,
		Files:   req.Files,
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
