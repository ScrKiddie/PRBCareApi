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

type KontrolBalikService struct {
	DB                     *gorm.DB
	KontrolBalikRepository *repository.KontrolBalikRepository
	PasienRepository       *repository.PasienRepository
	Validator              *validator.Validate
}

func NewKontrolBalikService(
	db *gorm.DB,
	kontrolBalikRepository *repository.KontrolBalikRepository,
	pasienRepository *repository.PasienRepository,
	validator *validator.Validate,
) *KontrolBalikService {
	return &KontrolBalikService{db, kontrolBalikRepository, pasienRepository, validator}
}

func (s *KontrolBalikService) Search(ctx context.Context, request *model.KontrolBalikSearchRequest) (*[]model.KontrolBalikResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrBadRequest
	}

	kontrolBalik := new([]entity.KontrolBalik)
	if request.IdPengguna > 0 {
		if err := s.KontrolBalikRepository.SearchAsPengguna(tx, kontrolBalik, request.IdPengguna, request.Status); err != nil {
			slog.Error(err.Error())
			return nil, fiber.ErrInternalServerError
		}
	} else if request.IdAdminPuskesmas > 0 {
		if err := s.KontrolBalikRepository.SearchAsAdminPuskesmas(tx, kontrolBalik, request.IdAdminPuskesmas, request.Status); err != nil {
			slog.Error(err.Error())
			return nil, fiber.ErrInternalServerError
		}
	} else {
		if err := s.KontrolBalikRepository.Search(tx, kontrolBalik, request.Status); err != nil {
			slog.Error(err.Error())
			return nil, fiber.ErrInternalServerError
		}
	}

	var response []model.KontrolBalikResponse
	for _, k := range *kontrolBalik {
		response = append(response, model.KontrolBalikResponse{
			ID: k.ID,
			PasienResponse: &model.PasienResponse{
				ID:           k.Pasien.ID,
				NoRekamMedis: k.Pasien.NoRekamMedis,
				Pengguna: &model.PenggunaResponse{
					ID:              k.Pasien.Pengguna.ID,
					NamaLengkap:     k.Pasien.Pengguna.NamaLengkap,
					Telepon:         k.Pasien.Pengguna.Telepon,
					TeleponKeluarga: k.Pasien.Pengguna.TeleponKeluarga,
					Alamat:          k.Pasien.Pengguna.Alamat,
				},
				AdminPuskesmas: &model.AdminPuskesmasResponse{
					ID:               k.Pasien.AdminPuskesmas.ID,
					NamaPuskesmas:    k.Pasien.AdminPuskesmas.NamaPuskesmas,
					Telepon:          k.Pasien.AdminPuskesmas.Telepon,
					Alamat:           k.Pasien.AdminPuskesmas.Alamat,
					WaktuOperasional: k.Pasien.AdminPuskesmas.WaktuOperasional,
				},
				TanggalDaftar: k.Pasien.TanggalDaftar,
				Status:        k.Pasien.Status,
			},
			Keluhan:        k.Keluhan,
			BeratBadan:     k.BeratBadan,
			TinggiBadan:    k.TinggiBadan,
			TekananDarah:   k.TekananDarah,
			DenyutNadi:     k.DenyutNadi,
			HasilLab:       k.HasilLab,
			HasilEkg:       k.HasilEkg,
			HasilDiagnosa:  k.HasilDiagnosa,
			TanggalKontrol: k.TanggalKontrol,
			Status:         k.Status,
		})
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &response, nil
}

func (s *KontrolBalikService) Get(ctx context.Context, request *model.KontrolBalikGetRequest) (*model.KontrolBalikResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrBadRequest
	}

	kontrolBalik := new(entity.KontrolBalik)
	if request.IdAdminPuskesmas > 0 {
		if err := s.KontrolBalikRepository.FindByIdAndIdAdminPuskesmasAndStatusNot(tx, kontrolBalik, request.ID, request.IdAdminPuskesmas, constant.StatusKontrolBalikBatal, true, false); err != nil {
			slog.Error(err.Error())
			return nil, fiber.ErrNotFound
		}
	} else if err := s.KontrolBalikRepository.FindByIdAndStatusNot(tx, kontrolBalik, request.ID, constant.StatusKontrolBalikBatal, true, true); err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	response := new(model.KontrolBalikResponse)
	response.ID = kontrolBalik.ID
	response.Keluhan = kontrolBalik.Keluhan
	response.BeratBadan = kontrolBalik.BeratBadan
	response.TinggiBadan = kontrolBalik.TinggiBadan
	response.TekananDarah = kontrolBalik.TekananDarah
	response.DenyutNadi = kontrolBalik.DenyutNadi
	response.HasilLab = kontrolBalik.HasilLab
	response.HasilEkg = kontrolBalik.HasilEkg
	response.HasilDiagnosa = kontrolBalik.HasilDiagnosa
	response.TanggalKontrol = kontrolBalik.TanggalKontrol
	response.IdPasien = kontrolBalik.IdPasien
	response.PasienResponse = &model.PasienResponse{
		ID:           kontrolBalik.Pasien.ID,
		NoRekamMedis: kontrolBalik.Pasien.NoRekamMedis,
		Pengguna: &model.PenggunaResponse{
			NamaLengkap:     kontrolBalik.Pasien.Pengguna.NamaLengkap,
			Telepon:         kontrolBalik.Pasien.Pengguna.Telepon,
			TeleponKeluarga: kontrolBalik.Pasien.Pengguna.TeleponKeluarga,
			Alamat:          kontrolBalik.Pasien.Pengguna.Alamat,
		},
		AdminPuskesmas: &model.AdminPuskesmasResponse{
			ID:               kontrolBalik.Pasien.AdminPuskesmas.ID,
			NamaPuskesmas:    kontrolBalik.Pasien.AdminPuskesmas.NamaPuskesmas,
			Telepon:          kontrolBalik.Pasien.AdminPuskesmas.Telepon,
			Alamat:           kontrolBalik.Pasien.AdminPuskesmas.Alamat,
			WaktuOperasional: kontrolBalik.Pasien.AdminPuskesmas.WaktuOperasional,
		},
		TanggalDaftar: kontrolBalik.Pasien.TanggalDaftar,
		Status:        kontrolBalik.Pasien.Status,
	}

	return response, nil
}

func (s *KontrolBalikService) Create(ctx context.Context, request *model.KontrolBalikCreateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	pasien := new(entity.Pasien)
	if request.IdAdminPuskesmas > 0 {
		if err := s.PasienRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, pasien, request.IdPasien, request.IdAdminPuskesmas, constant.StatusPasienAktif); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.PasienRepository.FindByIdAndStatus(tx, pasien, request.IdPasien, constant.StatusPasienAktif); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	}

	kontrolBalik := new(entity.KontrolBalik)
	kontrolBalik.IdPasien = request.IdPasien
	kontrolBalik.TanggalKontrol = request.TanggalKontrol
	kontrolBalik.Status = constant.StatusKontrolBalikMenunggu

	if err := s.KontrolBalikRepository.Create(tx, kontrolBalik); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *KontrolBalikService) Update(ctx context.Context, request *model.KontrolBalikUpdateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	kontrolBalik := new(entity.KontrolBalik)
	if request.IdAdminPuskesmas > 0 {
		if err := s.KontrolBalikRepository.FindByIdAndIdAdminPuskesmasAndStatusNot(tx, kontrolBalik, request.ID, request.IdAdminPuskesmas, constant.StatusKontrolBalikBatal, false, false); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.KontrolBalikRepository.FindByIdAndStatusNot(tx, kontrolBalik, request.ID, constant.StatusKontrolBalikBatal, false, false); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	}
	pasien := new(entity.Pasien)
	if request.IdPasien == kontrolBalik.IdPasien {
		if request.IdAdminPuskesmas > 0 {
			if err := s.PasienRepository.FindByIdAndIdAdminPuskesmas(tx, pasien, request.IdPasien, request.IdAdminPuskesmas); err != nil {
				slog.Error(err.Error())
				return fiber.ErrNotFound
			}
		} else {
			if err := s.PasienRepository.FindById(tx, pasien, request.IdPasien); err != nil {
				slog.Error(err.Error())
				return fiber.ErrNotFound
			}
		}
	} else {
		if request.IdAdminPuskesmas > 0 {
			if err := s.PasienRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, pasien, request.IdPasien, request.IdAdminPuskesmas, constant.StatusPasienAktif); err != nil {
				slog.Error(err.Error())
				return fiber.ErrNotFound
			}
		} else {
			if err := s.PasienRepository.FindByIdAndStatus(tx, pasien, request.IdPasien, constant.StatusPasienAktif); err != nil {
				slog.Error(err.Error())
				return fiber.ErrNotFound
			}
		}
	}

	kontrolBalik.IdPasien = request.IdPasien
	kontrolBalik.Keluhan = request.Keluhan
	kontrolBalik.BeratBadan = request.BeratBadan
	kontrolBalik.TinggiBadan = request.TinggiBadan
	kontrolBalik.TekananDarah = request.TekananDarah
	kontrolBalik.DenyutNadi = request.DenyutNadi
	kontrolBalik.HasilLab = request.HasilLab
	kontrolBalik.HasilEkg = request.HasilEkg
	kontrolBalik.HasilDiagnosa = request.HasilDiagnosa
	kontrolBalik.TanggalKontrol = request.TanggalKontrol

	if err := s.KontrolBalikRepository.Update(tx, kontrolBalik); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *KontrolBalikService) Delete(ctx context.Context, request *model.KontrolBalikDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	kontrolBalik := new(entity.KontrolBalik)
	if request.IdAdminPuskesmas > 0 {
		if err := s.KontrolBalikRepository.FindByIdAndIdAdminPuskesmasAndStatusOrStatus(tx, kontrolBalik, request.ID, request.IdAdminPuskesmas, constant.StatusKontrolBalikBatal, constant.StatusKontrolBalikSelesai); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.KontrolBalikRepository.FindByIdAndStatusOrStatus(tx, kontrolBalik, request.ID, constant.StatusKontrolBalikBatal, constant.StatusKontrolBalikSelesai); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	}

	if err := s.KontrolBalikRepository.Delete(tx, kontrolBalik); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *KontrolBalikService) Batal(ctx context.Context, request *model.KontrolBalikBatalRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	kontrolBalik := new(entity.KontrolBalik)
	if request.IdAdminPuskesmas > 0 {
		if err := s.KontrolBalikRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, kontrolBalik, request.ID, request.IdAdminPuskesmas, constant.StatusKontrolBalikMenunggu); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.KontrolBalikRepository.FindByIdAndStatus(tx, kontrolBalik, request.ID, constant.StatusKontrolBalikMenunggu); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	}

	kontrolBalik.Status = constant.StatusKontrolBalikBatal

	if err := s.KontrolBalikRepository.Update(tx, kontrolBalik); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *KontrolBalikService) Selesai(ctx context.Context, request *model.KontrolBalikSelesaiRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		slog.Error(err.Error())
		return fiber.ErrBadRequest
	}

	kontrolBalik := new(entity.KontrolBalik)
	if request.IdAdminPuskesmas > 0 {
		if err := s.KontrolBalikRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, kontrolBalik, request.ID, request.IdAdminPuskesmas, constant.StatusKontrolBalikMenunggu); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.KontrolBalikRepository.FindByIdAndStatus(tx, kontrolBalik, request.ID, constant.StatusKontrolBalikMenunggu); err != nil {
			slog.Error(err.Error())
			return fiber.ErrNotFound
		}
	}

	kontrolBalik.Status = constant.StatusKontrolBalikSelesai

	if err := s.KontrolBalikRepository.Update(tx, kontrolBalik); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}
