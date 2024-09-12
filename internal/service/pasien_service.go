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

type PasienService struct {
	DB                        *gorm.DB
	PasienRepository          *repository.PasienRepository
	AdminPuskesmasRepository  *repository.AdminPuskesmasRepository
	PenggunaRepository        *repository.PenggunaRepository
	KontrolBalikRepository    *repository.KontrolBalikRepository
	PengambilanObatRepository *repository.PengambilanObatRepository
	Validator                 *validator.Validate
}

func NewPasienService(
	db *gorm.DB,
	pasienRepository *repository.PasienRepository,
	adminPuskesmasRepository *repository.AdminPuskesmasRepository,
	penggunaRepository *repository.PenggunaRepository,
	kontrolBalikRepository *repository.KontrolBalikRepository,
	pengambilanObatRepository *repository.PengambilanObatRepository,
	validator *validator.Validate,
) *PasienService {
	return &PasienService{db, pasienRepository, adminPuskesmasRepository, penggunaRepository, kontrolBalikRepository, pengambilanObatRepository, validator}
}

func (s *PasienService) Search(ctx context.Context, request *model.PasienSearchRequest) (*[]model.PasienResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrBadRequest
	}

	pasien := new([]entity.Pasien)
	if request.IdPengguna > 0 {
		if err := s.PasienRepository.SearchAsPengguna(tx, pasien, request.IdPengguna, request.Status); err != nil {
			slog.Error(err.Error())
			return nil, fiber.ErrInternalServerError
		}
	} else if request.IdAdminPuskesmas > 0 {
		if err := s.PasienRepository.SearchAsAdminPuskesmas(tx, pasien, request.IdAdminPuskesmas, request.Status); err != nil {
			slog.Error(err.Error())
			return nil, fiber.ErrInternalServerError
		}
	} else {
		if err := s.PasienRepository.Search(tx, pasien, request.Status); err != nil {
			slog.Error(err.Error())
			return nil, fiber.ErrInternalServerError
		}
	}

	var response []model.PasienResponse
	for _, p := range *pasien {
		response = append(response, model.PasienResponse{
			ID:            p.ID,
			NoRekamMedis:  p.NoRekamMedis,
			TanggalDaftar: p.TanggalDaftar,
			Status:        p.Status,
			Pengguna: &model.PenggunaResponse{
				ID:              p.Pengguna.ID,
				NamaLengkap:     p.Pengguna.NamaLengkap,
				Telepon:         p.Pengguna.Telepon,
				TeleponKeluarga: p.Pengguna.TeleponKeluarga,
				Alamat:          p.Pengguna.Alamat,
			},
			AdminPuskesmas: &model.AdminPuskesmasResponse{
				ID:               p.AdminPuskesmas.ID,
				NamaPuskesmas:    p.AdminPuskesmas.NamaPuskesmas,
				Telepon:          p.AdminPuskesmas.Telepon,
				Alamat:           p.AdminPuskesmas.Alamat,
				WaktuOperasional: p.AdminPuskesmas.WaktuOperasional,
			},
		})
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &response, nil
}

func (s *PasienService) Get(ctx context.Context, request *model.PasienGetRequest) (*model.PasienResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrBadRequest
	}

	pasien := new(entity.Pasien)
	if request.IdAdminPuskesmas > 0 {
		if err := s.PasienRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, pasien, request.ID, request.IdAdminPuskesmas, constant.StatusPasienAktif); err != nil {
			slog.Error(err.Error())
			return nil, fiber.ErrNotFound
		}
	} else if err := s.PasienRepository.FindByIdAndStatus(tx, pasien, request.ID, constant.StatusPasienAktif); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	response := new(model.PasienResponse)
	response.ID = pasien.ID
	response.NoRekamMedis = pasien.NoRekamMedis
	response.TanggalDaftar = pasien.TanggalDaftar
	response.IdAdminPuskesmas = pasien.IdAdminPuskesmas
	response.IdPengguna = pasien.IdPengguna

	return response, nil
}

func (s *PasienService) Create(ctx context.Context, request *model.PasienCreateRequest) error {
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

	if err := s.PenggunaRepository.FindById(tx, &entity.Pengguna{}, request.IdPengguna); err != nil {
		slog.Error(err.Error())
		return fiber.ErrNotFound
	}

	pasien := new(entity.Pasien)
	pasien.NoRekamMedis = request.NoRekamMedis
	pasien.IdPengguna = request.IdPengguna
	pasien.IdAdminPuskesmas = request.IdAdminPuskesmas
	pasien.TanggalDaftar = request.TanggalDaftar
	pasien.Status = constant.StatusPasienAktif

	if err := s.PasienRepository.Create(tx, pasien); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *PasienService) Update(ctx context.Context, request *model.PasienUpdateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	pasien := new(entity.Pasien)
	if request.CurrentAdminPuskesmas {
		if err := s.PasienRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, pasien, request.ID, request.IdAdminPuskesmas, constant.StatusPasienAktif); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	} else if err := s.PasienRepository.FindByIdAndStatus(tx, pasien, request.ID, constant.StatusPasienAktif); err != nil {
		slog.Error(err.Error())
		return fiber.ErrNotFound
	}

	if err := s.AdminPuskesmasRepository.FindById(tx, &entity.AdminPuskesmas{}, request.IdAdminPuskesmas); err != nil {
		slog.Error(err.Error())
		return fiber.ErrNotFound
	}
	if err := s.PenggunaRepository.FindById(tx, &entity.Pengguna{}, request.IdPengguna); err != nil {
		slog.Error(err.Error())
		return fiber.ErrNotFound
	}

	pasien.NoRekamMedis = request.NoRekamMedis
	pasien.IdPengguna = request.IdPengguna
	pasien.IdAdminPuskesmas = request.IdAdminPuskesmas
	pasien.TanggalDaftar = request.TanggalDaftar

	if err := s.PasienRepository.Update(tx, pasien); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *PasienService) Selesai(ctx context.Context, request *model.PasienSelesaiRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}
	pasien := new(entity.Pasien)
	if request.IdAdminPuskesmas > 0 {
		if err := s.PasienRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, pasien, request.ID, request.IdAdminPuskesmas, constant.StatusPasienAktif); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	} else if err := s.PasienRepository.FindByIdAndStatus(tx, pasien, request.ID, constant.StatusPasienAktif); err != nil {
		slog.Error(err.Error())
		return fiber.ErrNotFound
	}

	if err := s.KontrolBalikRepository.FindByIdPasienAndStatus(tx, &entity.KontrolBalik{}, request.ID, constant.StatusKontrolBalikMenunggu); err == nil {
		return fiber.NewError(fiber.StatusConflict, "Pasien masih memiliki kontrol balik yang harus dilakukan")
	}

	if err := s.PengambilanObatRepository.FindByIdPasienAndStatus(tx, &entity.PengambilanObat{}, request.ID, constant.StatusKontrolBalikMenunggu); err == nil {
		return fiber.NewError(fiber.StatusConflict, "Pasien masih memiliki pengambilan obat yang harus dilakukan")
	}

	pasien.Status = constant.StatusPasienSelesai

	if err := s.PasienRepository.Update(tx, pasien); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil

}

func (s *PasienService) Delete(ctx context.Context, request *model.PasienDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	pasien := new(entity.Pasien)
	if request.IdAdminPuskesmas > 0 {
		if err := s.PasienRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, pasien, request.ID, request.IdAdminPuskesmas, constant.StatusPasienSelesai); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	} else if err := s.PasienRepository.FindByIdAndStatus(tx, pasien, request.ID, constant.StatusPasienSelesai); err != nil {
		slog.Error(err.Error())
		return fiber.ErrNotFound
	}

	if err := s.KontrolBalikRepository.FindByIdPasien(tx, &entity.KontrolBalik{}, request.ID); err == nil {
		return fiber.NewError(fiber.StatusConflict, "Pasien masih terkait dengan data kontrol balik yang ada")
	}
	if err := s.PengambilanObatRepository.FindByIdPasien(tx, &entity.PengambilanObat{}, request.ID); err == nil {
		return fiber.NewError(fiber.StatusConflict, "Pasien masih terkait dengan data pengambilan obat yang ada")
	}

	if err := s.PasienRepository.Delete(tx, pasien); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}
