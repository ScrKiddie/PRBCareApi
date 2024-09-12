package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"log/slog"
	"prb_care_api/internal/constant"
	"prb_care_api/internal/model"
	"prb_care_api/internal/service"
	"strings"
)

func AuthMiddleware(config *viper.Viper, adminSuperService *service.AdminSuperService, adminPuskesmasService *service.AdminPuskesmasService, adminApotekService *service.AdminApotekService, penggunaService *service.PenggunaService) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		tokenWithBearer := ctx.Get("Authorization")
		if tokenWithBearer == "" {
			return fiber.ErrUnauthorized
		}

		parts := strings.Split(tokenWithBearer, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.ErrUnauthorized
		}
		token := parts[1]

		tokenVerify, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				slog.Error("Unexpected", "signing method", token.Header["alg"])
				return nil, fiber.ErrInternalServerError
			}
			return []byte(config.GetString("jwt.secret")), nil
		})
		if err != nil {
			slog.Error(err.Error())
			return fiber.ErrUnauthorized
		}

		var id int32
		var role string
		if claims, ok := tokenVerify.Claims.(jwt.MapClaims); ok && tokenVerify.Valid {
			if subFloat64, ok := claims["sub"].(float64); ok {
				id = int32(subFloat64)
			} else {
				return fiber.ErrUnauthorized
			}
			if roleString, ok := claims["role"].(string); ok {
				role = roleString
			} else {
				return fiber.ErrUnauthorized
			}
		} else {
			return fiber.ErrUnauthorized
		}

		if role == constant.RoleAdminSuper {
			request := &model.AdminSuperVerifyRequest{ID: id}
			if err := adminSuperService.Verify(ctx.UserContext(), request); err == nil {
				slog.Info("Authenticated as", "AdminSuper:", id)
				auth := &model.Auth{ID: id, Role: role}
				ctx.Locals("auth", auth)
				return ctx.Next()
			}
		} else if role == constant.RoleAdminPuskesmas {
			request := &model.AdminPuskesmasVerifyRequest{ID: id}
			if err := adminPuskesmasService.Verify(ctx.UserContext(), request); err == nil {
				slog.Info("Authenticated as", "AdminPuskesmas", id)
				auth := &model.Auth{ID: id, Role: role}
				ctx.Locals("auth", auth)
				return ctx.Next()
			}
		} else if role == constant.RoleAdminApotek {
			request := &model.AdminApotekVerifyRequest{ID: id}
			if err := adminApotekService.Verify(ctx.UserContext(), request); err == nil {
				slog.Info("Authenticated as", "AdminApotek:", id)
				auth := &model.Auth{ID: id, Role: role}
				ctx.Locals("auth", auth)
				return ctx.Next()
			}
		} else if role == constant.RolePengguna {
			request := &model.PenggunaVerifyRequest{ID: id}
			if err := penggunaService.Verify(ctx.UserContext(), request); err == nil {
				slog.Info("Authenticated as", "Pengguna", id)
				auth := &model.Auth{ID: id, Role: role}
				ctx.Locals("auth", auth)
				return ctx.Next()
			}
		}
		return fiber.ErrUnauthorized
	}
}
func GetAuth(ctx fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
