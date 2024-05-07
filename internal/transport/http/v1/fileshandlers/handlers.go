package fileshandlers

import (
	"backend/internal/services"
	"fmt"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

type handler struct {
	service services.FileService
	log     *zerolog.Logger
}

func (h *handler) getById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, `Path parameter <id> empty or not a number`)
	}

	file, err := h.service.GetById(ctx.UserContext(), int64(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.GetById: %v", err))
	}

	responseBytes, err := jsoniter.Marshal(file)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("jsoniter.Marshal: %v", err))
	}

	if err = ctx.Status(fiber.StatusOK).Send(responseBytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.Send: %v", err))
	}

	return nil
}

func (h *handler) create(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("ctx.MultipartForm: %v", err))
	}

	var req createRequest
	if err = jsoniter.UnmarshalFromString(form.Value["body"][0], &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("jsoniter.UnmarshalFromString: %v", err))
	}

	for _, header := range form.File["file"] {
		originalName := header.Filename

		saveName := uuid.Must(uuid.NewV7()).String() + filepath.Ext(header.Filename)
		savePath := "./storage"

		if err = ctx.SaveFile(header, savePath+"/"+saveName); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.SaveFile: %v", err))
		}

		_, err = h.service.Create(ctx.UserContext(), services.FileServiceCreateOpts{
			Name:     originalName,
			Filename: saveName,
			Filepath: savePath,
			TaskId:   req.TaskId,
			AnswerId: req.AnswerId,
		})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.Create: %v", err))
		}
	}

	if err = ctx.Status(fiber.StatusCreated).Send(nil); err != nil {
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
