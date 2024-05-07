package authhandlers

import (
	"backend/internal/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

const refreshTokenCookieName = "refreshToken"

type handler struct {
	service services.UserService
	log     *zerolog.Logger
}

func (h *handler) signUp(ctx *fiber.Ctx) error {
	var request signUpRequest
	if err := jsoniter.Unmarshal(ctx.Body(), &request); err != nil {
		h.log.Error().Err(err).Send()
		return fiber.NewError(fiber.StatusBadRequest, "SignUp request not valid")
	}

	_, err := h.service.Create(ctx.UserContext(), services.UserServiceCreateOpts{
		GroupId:    request.GroupId,
		Email:      request.Email,
		Password:   request.Password,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		MiddleName: request.MiddleName,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.service.Create: %v", err))
	}

	if err = ctx.Status(fiber.StatusCreated).Send(nil); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("ctx.Send: %v", err))
	}

	return nil
}

//func (h *handler) signIn(ctx *fiber.Ctx) error {
//	var request signInRequest
//	if err := jsoniter.Unmarshal(ctx.Body(), &request); err != nil {
//		h.log.Error().Err(err).Send()
//		return fiber.NewError(fiber.StatusBadRequest, "SignIn request not valid")
//	}
//
//	jwtPair, err := h.authService.SignIn(ctx.UserContext(), entity.Credentials{
//		Login:    request.Login,
//		Password: request.Password,
//	})
//	if err != nil {
//		h.log.Error().Err(err).Send()
//		if errors.Is(err, auth.ErrNoAuth) {
//			return fiber.NewError(fiber.StatusUnauthorized, "Incorrect login or password")
//		}
//		return ctx.SendStatus(fiber.StatusInternalServerError)
//	}
//
//	cookie := new(fiber.Cookie)
//	cookie.Name = refreshTokenCookieName
//	cookie.Value = jwtPair.RefreshToken
//	cookie.HTTPOnly = true
//	cookie.Expires = time.Now().Add(h.authService.RefreshTokenExpTime())
//
//	ctx.Cookie(cookie)
//
//	responseBytes, err := jsoniter.Marshal(jwtPair)
//	if err != nil {
//		h.log.Error().Err(err).Send()
//		return ctx.SendStatus(fiber.StatusInternalServerError)
//	}
//
//	if err = ctx.Status(fiber.StatusOK).Send(responseBytes); err != nil {
//		h.log.Error().Err(err).Send()
//		return ctx.SendStatus(fiber.StatusInternalServerError)
//	}
//
//	return nil
//}
//
//func (h *handler) refresh(ctx *fiber.Ctx) error {
//	refreshToken := ctx.Cookies("refreshToken")
//
//	jwtPair, err := h.authService.RefreshTokens(ctx.UserContext(), refreshToken)
//	if err != nil {
//		h.log.Error().Err(err).Send()
//		return ctx.SendStatus(fiber.StatusInternalServerError)
//	}
//
//	cookie := new(fiber.Cookie)
//	cookie.Name = refreshTokenCookieName
//	cookie.Value = jwtPair.RefreshToken
//	cookie.HTTPOnly = true
//	cookie.Expires = time.Now().Add(h.authService.RefreshTokenExpTime())
//
//	ctx.Cookie(cookie)
//
//	responseBytes, err := jsoniter.Marshal(jwtResponse{
//		AccessToken: jwtPair.AccessToken,
//	})
//	if err != nil {
//		h.log.Error().Err(err).Send()
//		return ctx.SendStatus(fiber.StatusInternalServerError)
//	}
//
//	if err = ctx.Status(fiber.StatusOK).Send(responseBytes); err != nil {
//		h.log.Error().Err(err).Send()
//		return ctx.SendStatus(fiber.StatusInternalServerError)
//	}
//
//	return nil
//}
