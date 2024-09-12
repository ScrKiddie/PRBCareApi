package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"log/slog"
	"prb_care_api/internal/entity"
	"prb_care_api/internal/model"
	"prb_care_api/internal/repository"
	"time"
)

type ArtikelService struct {
	DB                       *gorm.DB
	ArtikelRepository        *repository.ArtikelRepository
	AdminPuskesmasRepository *repository.AdminPuskesmasRepository
	Validator                *validator.Validate
}

func NewArtikelService(
	db *gorm.DB,
	artikelRepository *repository.ArtikelRepository,
	adminPuskesmasRepository *repository.AdminPuskesmasRepository,
	validator *validator.Validate,
) *ArtikelService {
	return &ArtikelService{
		DB:                       db,
		ArtikelRepository:        artikelRepository,
		AdminPuskesmasRepository: adminPuskesmasRepository,
		Validator:                validator,
	}
}

func (s *ArtikelService) Search(ctx context.Context, request *model.ArtikelSearchRequest) (*[]model.ArtikelResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrBadRequest
	}

	artikel := new([]entity.Artikel)
	if err := s.ArtikelRepository.Search(tx, artikel, request.IdAdminPuskesmas); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	var response []model.ArtikelResponse
	for _, a := range *artikel {
		response = append(response, model.ArtikelResponse{
			ID:               a.ID,
			Judul:            a.Judul,
			Ringkasan:        a.Ringkasan,
			TanggalPublikasi: a.TanggalPublikasi,
			AdminPuskesmas: &model.AdminPuskesmasResponse{
				ID:               a.AdminPuskesmas.ID,
				NamaPuskesmas:    a.AdminPuskesmas.NamaPuskesmas,
				Telepon:          a.AdminPuskesmas.Telepon,
				Alamat:           a.AdminPuskesmas.Alamat,
				WaktuOperasional: a.AdminPuskesmas.WaktuOperasional,
			},
		})
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &response, nil
}

func (s *ArtikelService) Get(ctx context.Context, request *model.ArtikelGetRequest) (*model.ArtikelResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrBadRequest
	}

	artikel := new(entity.Artikel)
	if err := s.ArtikelRepository.FindById(tx, artikel, request.ID); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrNotFound
	}

	response := &model.ArtikelResponse{
		ID:               artikel.ID,
		Judul:            artikel.Judul,
		Ringkasan:        artikel.Ringkasan,
		Isi:              artikel.Isi,
		TanggalPublikasi: artikel.TanggalPublikasi,
		AdminPuskesmas: &model.AdminPuskesmasResponse{
			ID:               artikel.AdminPuskesmas.ID,
			NamaPuskesmas:    artikel.AdminPuskesmas.NamaPuskesmas,
			Telepon:          artikel.AdminPuskesmas.Telepon,
			Alamat:           artikel.AdminPuskesmas.Alamat,
			WaktuOperasional: artikel.AdminPuskesmas.WaktuOperasional,
		},
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return response, nil
}

func (s *ArtikelService) Create(ctx context.Context, request *model.ArtikelCreateRequest) error {
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

	artikel := &entity.Artikel{
		Judul:            request.Judul,
		Ringkasan:        request.Ringkasan,
		Isi:              request.Isi,
		TanggalPublikasi: time.Now().Unix(),
		IdAdminPuskesmas: request.IdAdminPuskesmas,
	}

	if err := s.ArtikelRepository.Create(tx, artikel); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *ArtikelService) Update(ctx context.Context, request *model.ArtikelUpdateRequest) error {
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

	artikel := new(entity.Artikel)
	if request.CurrentAdminPuskesmas {
		if err := s.ArtikelRepository.FindByIdAndIdAdminPuskesmas(tx, artikel, request.IdAdminPuskesmas, request.ID); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.ArtikelRepository.FindById(tx, artikel, request.ID); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	}

	artikel.Judul = request.Judul
	artikel.Ringkasan = request.Ringkasan
	artikel.Isi = request.Isi
	artikel.IdAdminPuskesmas = request.IdAdminPuskesmas

	if err := s.ArtikelRepository.Update(tx, artikel); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *ArtikelService) Delete(ctx context.Context, request *model.ArtikelDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	artikel := new(entity.Artikel)
	if request.IdAdminPuskesmas > 0 {
		if err := s.ArtikelRepository.FindByIdAndIdAdminPuskesmas(tx, artikel, request.IdAdminPuskesmas, request.ID); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.ArtikelRepository.FindById(tx, artikel, request.ID); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	}

	if err := s.ArtikelRepository.Delete(tx, artikel); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}
