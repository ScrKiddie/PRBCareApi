package service

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"prb_care_api/internal/constant"
	"prb_care_api/internal/entity"
	"prb_care_api/internal/model"
	"prb_care_api/internal/repository"
	"time"
)

type PenggunaService struct {
	DB                 *gorm.DB
	PenggunaRepository *repository.PenggunaRepository
	PasienRepository   *repository.PasienRepository
	Validator          *validator.Validate
	Config             *viper.Viper
}

func NewPenggunaService(db *gorm.DB,
	penggunaRepository *repository.PenggunaRepository,
	pasienRepository *repository.PasienRepository,
	validator *validator.Validate,
	config *viper.Viper) *PenggunaService {
	return &PenggunaService{db, penggunaRepository, pasienRepository, validator, config}
}

func (s *PenggunaService) List(ctx context.Context) (*[]model.PenggunaResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	pengguna := new([]entity.Pengguna)
	if err := s.PenggunaRepository.FindAll(tx, pengguna); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	var response []model.PenggunaResponse
	for _, p := range *pengguna {
		response = append(response, model.PenggunaResponse{
			ID:              p.ID,
			NamaLengkap:     p.NamaLengkap,
			Telepon:         p.Telepon,
			TeleponKeluarga: p.TeleponKeluarga,
			Alamat:          p.Alamat,
		})
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &response, nil
}

func (s *PenggunaService) Get(ctx context.Context, request *model.PenggunaGetRequest) (*model.PenggunaResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	pengguna := new(entity.Pengguna)
	if err := s.PenggunaRepository.FindById(tx, pengguna, request.ID); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	response := new(model.PenggunaResponse)
	response.ID = pengguna.ID
	response.Username = pengguna.Username
	response.NamaLengkap = pengguna.NamaLengkap
	response.Alamat = pengguna.Alamat
	response.Telepon = pengguna.Telepon
	response.TeleponKeluarga = pengguna.TeleponKeluarga

	return response, nil
}

func (s *PenggunaService) Create(ctx context.Context, request *model.PenggunaCreateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	total, err := s.PenggunaRepository.CountByUsername(tx, request.Username)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	if total > 0 {
		return fiber.NewError(fiber.StatusConflict, "Username sudah digunakan")
	}

	total, err = s.PenggunaRepository.CountByTelepon(tx, request.Telepon)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	if total > 0 {
		return fiber.NewError(fiber.StatusConflict, "Telepon sudah digunakan")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	penggunaEnity := new(entity.Pengguna)
	penggunaEnity.Username = request.Username
	penggunaEnity.NamaLengkap = request.NamaLengkap
	penggunaEnity.Alamat = request.Alamat
	penggunaEnity.Telepon = request.Telepon
	penggunaEnity.TeleponKeluarga = request.TeleponKeluarga
	penggunaEnity.Password = string(password)

	if err := s.PenggunaRepository.Create(tx, penggunaEnity); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *PenggunaService) Update(ctx context.Context, request *model.PenggunaUpdateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	pengguna := new(entity.Pengguna)
	if err := s.PenggunaRepository.FindById(tx, pengguna, request.ID); err != nil {
		log.Println(err.Error())
		return fiber.ErrNotFound
	}

	total, err := s.PenggunaRepository.CountByUsername(tx, request.Username)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	if total > 0 && pengguna.Username != request.Username {
		return fiber.NewError(fiber.StatusConflict, "Username sudah digunakan")
	}

	total, err = s.PenggunaRepository.CountByTelepon(tx, request.Telepon)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	if total > 0 && pengguna.Telepon != request.Telepon {
		return fiber.NewError(fiber.StatusConflict, "Telepon sudah digunakan")
	}

	var password []byte
	if request.Password != "" {
		password, err = bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err.Error())
			return fiber.ErrInternalServerError
		}
	}

	pengguna.Username = request.Username
	pengguna.NamaLengkap = request.NamaLengkap
	pengguna.Alamat = request.Alamat
	pengguna.Telepon = request.Telepon
	pengguna.TeleponKeluarga = request.TeleponKeluarga
	if string(password) != "" {
		pengguna.Password = string(password)
	}

	if err := s.PenggunaRepository.Update(tx, pengguna); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *PenggunaService) Delete(ctx context.Context, request *model.PenggunaDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	pengguna := new(entity.Pengguna)
	if err := s.PenggunaRepository.FindById(tx, pengguna, request.ID); err != nil {
		log.Println(err.Error())
		return fiber.ErrNotFound
	}

	if err := s.PasienRepository.FindByIdPengguna(tx, &entity.Pasien{}, request.ID); err == nil {
		return fiber.NewError(fiber.StatusConflict, "Pengguna masih terkait dengan data pasien yang ada")
	}

	if err := s.PenggunaRepository.Delete(tx, pengguna); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *PenggunaService) Login(ctx context.Context, request *model.PenggunaLoginRequest) (*model.PenggunaResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	pengguna := new(entity.Pengguna)
	if err := s.PenggunaRepository.FindByUsername(tx, pengguna, request.Username); err != nil {
		log.Println(err.Error())
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Username atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pengguna.Password), []byte(request.Password)); err != nil {
		log.Println(err.Error())
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Username atau password salah")
	}

	key := s.Config.GetString("jwt.secret")
	exp := s.Config.GetInt("jwt.exp")
	claims := jwt.MapClaims{
		"sub":  pengguna.ID,
		"role": constant.RolePengguna,
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

	return &model.PenggunaResponse{Token: token}, nil
}

func (s *PenggunaService) Verify(ctx context.Context, request *model.PenggunaVerifyRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	pengguna := new(entity.Pengguna)
	if err := s.PenggunaRepository.FindById(tx, pengguna, request.ID); err != nil {
		log.Println(err.Error())
		return fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *PenggunaService) Current(ctx context.Context, request *model.PenggunaGetRequest) (*model.PenggunaResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	pengguna := new(entity.Pengguna)
	if err := s.PenggunaRepository.FindById(tx, pengguna, request.ID); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	response := new(model.PenggunaResponse)
	response.NamaLengkap = pengguna.NamaLengkap
	response.Alamat = pengguna.Alamat
	response.Telepon = pengguna.Telepon
	response.TeleponKeluarga = pengguna.TeleponKeluarga

	return response, nil
}

func (s *PenggunaService) CurrentProfileUpdate(ctx context.Context, request *model.PenggunaProfileUpdateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	pengguna := new(entity.Pengguna)
	if err := s.PenggunaRepository.FindById(tx, pengguna, request.ID); err != nil {
		log.Println(err.Error())
		return fiber.ErrNotFound
	}

	total, err := s.PenggunaRepository.CountByTelepon(tx, request.Telepon)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	fmt.Println(total)
	if total > 0 && pengguna.Telepon != request.Telepon {
		return fiber.NewError(fiber.StatusConflict, "Telepon sudah digunakan")
	}

	pengguna.NamaLengkap = request.NamaLengkap
	pengguna.Alamat = request.Alamat
	pengguna.Telepon = request.Telepon
	pengguna.TeleponKeluarga = request.TeleponKeluarga

	if err := s.PenggunaRepository.Update(tx, pengguna); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *PenggunaService) CurrentPasswordUpdate(ctx context.Context, request *model.PenggunaPasswordUpdateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	pengguna := new(entity.Pengguna)
	if err := s.PenggunaRepository.FindById(tx, pengguna, request.ID); err != nil {
		log.Println(err.Error())
		return fiber.ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pengguna.Password), []byte(request.CurrentPassword)); err != nil {
		log.Println(err.Error())
		return fiber.NewError(fiber.StatusUnauthorized, "Password saat ini salah")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}
	pengguna.Password = string(password)

	if err := s.PenggunaRepository.Update(tx, pengguna); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *PenggunaService) CurrentTokenPerangkatUpdate(ctx context.Context, request *model.PenggunaTokenPerangkatUpdateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	pengguna := new(entity.Pengguna)
	if err := s.PenggunaRepository.FindById(tx, pengguna, request.ID); err != nil {
		log.Println(err.Error())
		return fiber.ErrNotFound
	}

	pengguna.TokenPerangkat = request.TokenPerangkat

	if err := s.PenggunaRepository.Update(tx, pengguna); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}
