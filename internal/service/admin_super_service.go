package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log/slog"
	"prb_care_api/internal/adapter"
	"prb_care_api/internal/constant"
	"prb_care_api/internal/entity"
	"prb_care_api/internal/model"
	"prb_care_api/internal/repository"
	"time"
)

type AdminSuperService struct {
	DB                   *gorm.DB
	AdminSuperRepository *repository.AdminSuperRepository
	RecaptchaAdapter     *adapter.Captcha
	Validator            *validator.Validate
	Config               *viper.Viper
}

func NewAdminSuperService(db *gorm.DB,
	adminSuperRepository *repository.AdminSuperRepository,
	captchaAdapter *adapter.Captcha,
	validator *validator.Validate,
	config *viper.Viper) *AdminSuperService {
	return &AdminSuperService{db, adminSuperRepository, captchaAdapter, validator, config}
}

func (s *AdminSuperService) Login(ctx context.Context, request *model.AdminSuperLoginRequest) (*model.AdminSuperResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrBadRequest
	}

	captchaRequest := &model.CaptchaRequest{
		TokenCaptcha: request.TokenCaptcha,
		Secret:       s.Config.GetString("captcha.secret"),
	}

	ok, err := s.RecaptchaAdapter.Verify(captchaRequest)
	if err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrInternalServerError
	}
	if !ok {
		return nil, fiber.ErrBadRequest
	}

	adminSuper := new(entity.AdminSuper)
	if err := s.AdminSuperRepository.FindByUsername(tx, adminSuper, request.Username); err != nil {
		slog.Error(err.Error())
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Username atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(adminSuper.Password), []byte(request.Password)); err != nil {
		slog.Error(err.Error())
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Username atau password salah")
	}

	key := s.Config.GetString("jwt.secret")
	exp := s.Config.GetInt("jwt.exp")
	claims := jwt.MapClaims{
		"sub":  adminSuper.ID,
		"role": constant.RoleAdminSuper,
		"exp":  time.Now().Add(time.Duration(exp) * time.Hour).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(key))
	if err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &model.AdminSuperResponse{Token: token}, nil
}

func (s *AdminSuperService) PasswordUpdate(ctx context.Context, request *model.AdminSuperPasswordUpdateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	adminSuper := new(entity.AdminSuper)
	if err := s.AdminSuperRepository.FindById(tx, adminSuper, request.ID); err != nil {
		slog.Error(err.Error())
		return fiber.ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(adminSuper.Password), []byte(request.CurrentPassword)); err != nil {
		slog.Error(err.Error())
		return fiber.NewError(fiber.StatusUnauthorized, "Password saat ini salah")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}
	adminSuper.Password = string(password)

	if err := s.AdminSuperRepository.Update(tx, adminSuper); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *AdminSuperService) Verify(ctx context.Context, request *model.AdminSuperVerifyRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	adminSuper := new(entity.AdminSuper)
	if err := s.AdminSuperRepository.FindById(tx, adminSuper, request.ID); err != nil {
		slog.Error(err.Error())
		return fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}
