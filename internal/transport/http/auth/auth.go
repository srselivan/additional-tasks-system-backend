package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
)

type Claims struct {
	UserId  int64
	GroupId *int64
	Role    int
}

func GetClaimsFromCtx(ctx *fiber.Ctx) (Claims, error) {
	claims := ctx.Locals("claims")

	claimsMap, ok := claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, fiber.NewError(fiber.StatusUnauthorized, "missed claims in token")
	}

	result := Claims{}

	userId, ok := claimsMap["sub"].(float64)
	if !ok {
		return Claims{}, fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}
	result.UserId = int64(userId)

	groupId, ok := claimsMap["grp"]
	if !ok {
		return Claims{}, fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}
	if id, ok := groupId.(float64); ok {
		result.GroupId = lo.ToPtr(int64(id))
	} else {
		result.GroupId = nil
	}

	role, ok := claimsMap["role"].(float64)
	if !ok {
		return Claims{}, fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}
	result.Role = int(role)

	return result, nil
}
