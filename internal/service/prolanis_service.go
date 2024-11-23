package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"log/slog"
	"prb_care_api/internal/constant"
	"prb_care_api/internal/entity"
	"prb_care_api/internal/model"
	"prb_care_api/internal/repository"
)

type ProlanisService struct {
	DB                       *gorm.DB
	ProlanisRepository       *repository.ProlanisRepository
	AdminPuskesmasRepository *repository.AdminPuskesmasRepository
	Validator                *validator.Validate
}

func NewProlanisService(
	db *gorm.DB,
	prolanisRepository *repository.ProlanisRepository,
	adminPuskesmasRepository *repository.AdminPuskesmasRepository,
	validator *validator.Validate,
) *ProlanisService {
	return &ProlanisService{db, prolanisRepository, adminPuskesmasRepository, validator}
}

func (s *ProlanisService) Search(ctx context.Context, request *model.ProlanisSearchRequest) (*[]model.ProlanisResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrBadRequest
	}

	prolanis := new([]entity.Prolanis)
	if err := s.ProlanisRepository.Search(tx, prolanis, request.Status, request.IdAdminPuskesmas); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	var response []model.ProlanisResponse
	for _, p := range *prolanis {
		response = append(response, model.ProlanisResponse{
			ID: p.ID,
			AdminPuskesmas: &model.AdminPuskesmasResponse{
				ID:               p.AdminPuskesmas.ID,
				NamaPuskesmas:    p.AdminPuskesmas.NamaPuskesmas,
				Telepon:          p.AdminPuskesmas.Telepon,
				Alamat:           p.AdminPuskesmas.Alamat,
				WaktuOperasional: p.AdminPuskesmas.WaktuOperasional,
			},
			Deskripsi:    p.Deskripsi,
			WaktuMulai:   p.WaktuMulai,
			WaktuSelesai: p.WaktuSelesai,
			Status:       p.Status,
		})
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &response, nil
}

func (s *ProlanisService) Get(ctx context.Context, request *model.ProlanisGetRequest) (*model.ProlanisResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrBadRequest
	}

	prolanis := new(entity.Prolanis)
	if request.IdAdminPuskesmas > 0 {
		if err := s.ProlanisRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, prolanis, request.ID, request.IdAdminPuskesmas, constant.StatusProlanisAktif); err != nil {
			slog.Error(err.Error())
			return nil, fiber.ErrNotFound
		}
	} else if err := s.ProlanisRepository.FindByIdAndStatus(tx, prolanis, request.ID, constant.StatusProlanisAktif); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	response := new(model.ProlanisResponse)
	response.ID = prolanis.ID
	response.Deskripsi = prolanis.Deskripsi
	response.WaktuMulai = prolanis.WaktuMulai
	response.WaktuSelesai = prolanis.WaktuSelesai
	response.IdAdminPuskesmas = prolanis.IdAdminPuskesmas

	return response, nil
}

func (s *ProlanisService) Create(ctx context.Context, request *model.ProlanisCreateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	if err := s.AdminPuskesmasRepository.FindById(tx, &entity.AdminPuskesmas{}, request.IdAdminPuskesmas); err != nil {
		slog.Error(err.Error())
		return fiber.ErrNotFound
	}

	prolanis := new(entity.Prolanis)
	prolanis.Deskripsi = request.Deskripsi
	prolanis.WaktuMulai = request.WaktuMulai
	prolanis.WaktuSelesai = request.WaktuSelesai
	prolanis.IdAdminPuskesmas = request.IdAdminPuskesmas
	prolanis.Status = constant.StatusProlanisAktif

	if err := s.ProlanisRepository.Create(tx, prolanis); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *ProlanisService) Update(ctx context.Context, request *model.ProlanisUpdateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	prolanis := new(entity.Prolanis)
	if request.CurrentAdminPuskesmas {
		if err := s.ProlanisRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, prolanis, request.ID, request.IdAdminPuskesmas, constant.StatusProlanisAktif); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	} else if err := s.ProlanisRepository.FindByIdAndStatus(tx, prolanis, request.ID, constant.StatusProlanisAktif); err != nil {
		slog.Error(err.Error())
		return fiber.ErrNotFound
	}

	if err := s.AdminPuskesmasRepository.FindById(tx, &entity.AdminPuskesmas{}, request.IdAdminPuskesmas); err != nil {
		slog.Error(err.Error())
		return fiber.ErrNotFound
	}

	prolanis.Deskripsi = request.Deskripsi
	prolanis.WaktuMulai = request.WaktuMulai
	prolanis.WaktuSelesai = request.WaktuSelesai
	prolanis.IdAdminPuskesmas = request.IdAdminPuskesmas

	if err := s.ProlanisRepository.Update(tx, prolanis); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *ProlanisService) Selesai(ctx context.Context, request *model.ProlanisSelesaiRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}
	prolanis := new(entity.Prolanis)
	if request.IdAdminPuskesmas > 0 {
		if err := s.ProlanisRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, prolanis, request.ID, request.IdAdminPuskesmas, constant.StatusProlanisAktif); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	} else if err := s.ProlanisRepository.FindByIdAndStatus(tx, prolanis, request.ID, constant.StatusProlanisAktif); err != nil {
		slog.Error(err.Error())
		return fiber.ErrNotFound
	}

	prolanis.Status = constant.StatusProlanisSelesai

	if err := s.ProlanisRepository.Update(tx, prolanis); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil

}

func (s *ProlanisService) Delete(ctx context.Context, request *model.ProlanisDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	prolanis := new(entity.Prolanis)
	if request.IdAdminPuskesmas > 0 {
		if err := s.ProlanisRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, prolanis, request.ID, request.IdAdminPuskesmas, constant.StatusProlanisSelesai); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	} else if err := s.ProlanisRepository.FindByIdAndStatus(tx, prolanis, request.ID, constant.StatusProlanisSelesai); err != nil {
		slog.Error(err.Error())
		return fiber.ErrNotFound
	}

	if err := s.ProlanisRepository.Delete(tx, prolanis); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}
