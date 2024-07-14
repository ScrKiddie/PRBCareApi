package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"prbcare_be/internal/constant"
	"prbcare_be/internal/entity"
	"prbcare_be/internal/model"
	"prbcare_be/internal/repository"
	"time"
)

type AdminSuperService struct {
	DB                   *gorm.DB
	AdminSuperRepository *repository.AdminSuperRepository
	Validator            *validator.Validate
	Config               *viper.Viper
}

func NewAdminSuperService(db *gorm.DB,
	adminSuperRepository *repository.AdminSuperRepository,
	validator *validator.Validate,
	config *viper.Viper) *AdminSuperService {
	return &AdminSuperService{db, adminSuperRepository, validator, config}
}

func (s *AdminSuperService) Login(ctx context.Context, request *model.AdminSuperLoginRequest) (*model.AdminSuperResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	adminSuper := new(entity.AdminSuper)
	if err := s.AdminSuperRepository.FindByUsername(tx, adminSuper, request.Username); err != nil {
		log.Println(err.Error())
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Username atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(adminSuper.Password), []byte(request.Password)); err != nil {
		log.Println(err.Error())
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
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &model.AdminSuperResponse{Token: token}, nil
}

func (s *AdminSuperService) PasswordUpdate(ctx context.Context, request *model.AdminSuperPasswordUpdateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	adminSuper := new(entity.AdminSuper)
	if err := s.AdminSuperRepository.FindById(tx, adminSuper, request.ID); err != nil {
		log.Println(err.Error())
		return fiber.NewError(fiber.StatusNotFound, "Not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(adminSuper.Password), []byte(request.CurrentPassword)); err != nil {
		log.Println(err.Error())
		return fiber.NewError(fiber.StatusUnauthorized, "Password saat ini salah")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	adminSuper.Password = string(password)

	if err := s.AdminSuperRepository.Update(tx, adminSuper); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *AdminSuperService) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	tokenString := request.Token

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, fiber.ErrInternalServerError
		}
		return []byte(s.Config.GetString("jwt.secret")), nil
	})

	if err != nil {
		log.Println("error parsing token:", err.Error())
		return nil, fiber.ErrUnauthorized
	}

	var id int
	var role string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if subFloat64, ok := claims["sub"].(float64); ok {
			id = int(subFloat64)
		} else {
			return nil, fiber.ErrUnauthorized
		}
		if roleString, ok := claims["role"].(string); ok {
			role = roleString
		} else {
			return nil, fiber.ErrUnauthorized
		}
	} else {
		return nil, fiber.ErrUnauthorized
	}

	if role != constant.RoleAdminSuper {
		return nil, fiber.ErrUnauthorized
	}

	user := new(entity.AdminSuper)
	if err := s.AdminSuperRepository.FindById(tx, user, id); err != nil {
		log.Println(err.Error())
		return nil, fiber.NewError(fiber.StatusNotFound, "Not found")
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &model.Auth{ID: user.ID}, nil
}
