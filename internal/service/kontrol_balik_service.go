package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"log"
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
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	kontrolBalik := new([]entity.KontrolBalik)
	if request.IdPengguna > 0 {
		if err := s.KontrolBalikRepository.SearchAsPengguna(tx, kontrolBalik, request.IdPengguna, request.Status); err != nil {
			log.Println(err.Error())
			return nil, fiber.ErrInternalServerError
		}
	} else if request.IdAdminPuskesmas > 0 {
		if err := s.KontrolBalikRepository.SearchAsAdminPuskesmas(tx, kontrolBalik, request.IdAdminPuskesmas, request.Status); err != nil {
			log.Println(err.Error())
			return nil, fiber.ErrInternalServerError
		}
	} else {
		if err := s.KontrolBalikRepository.Search(tx, kontrolBalik, request.Status); err != nil {
			log.Println(err.Error())
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
					ID:            k.Pasien.AdminPuskesmas.ID,
					NamaPuskesmas: k.Pasien.AdminPuskesmas.NamaPuskesmas,
					Telepon:       k.Pasien.AdminPuskesmas.Telepon,
					Alamat:        k.Pasien.AdminPuskesmas.Alamat,
				},
				BeratBadan:     k.Pasien.BeratBadan,
				TinggiBadan:    k.Pasien.TinggiBadan,
				TekananDarah:   k.Pasien.TekananDarah,
				DenyutNadi:     k.Pasien.DenyutNadi,
				HasilLab:       k.Pasien.HasilLab,
				HasilEkg:       k.Pasien.HasilEkg,
				TanggalPeriksa: k.Pasien.TanggalPeriksa,
				Status:         k.Pasien.Status,
			},
			TanggalKontrol: k.TanggalKontrol,
			Status:         k.Status,
		})
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &response, nil
}

func (s *KontrolBalikService) Get(ctx context.Context, request *model.KontrolBalikGetRequest) (*model.KontrolBalikResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	kontrolBalik := new(entity.KontrolBalik)
	if request.IdAdminPuskesmas > 0 {
		if err := s.KontrolBalikRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, kontrolBalik, request.ID, request.IdAdminPuskesmas, constant.StatusKontrolBalikMenunggu); err != nil {
			log.Println(err.Error())
			return nil, fiber.ErrNotFound
		}
	} else if err := s.KontrolBalikRepository.FindByIdAndStatus(tx, kontrolBalik, request.ID, constant.StatusKontrolBalikMenunggu); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	response := new(model.KontrolBalikResponse)
	response.ID = kontrolBalik.ID
	response.TanggalKontrol = kontrolBalik.TanggalKontrol
	response.IdPasien = kontrolBalik.IdPasien
	return response, nil
}

func (s *KontrolBalikService) Create(ctx context.Context, request *model.KontrolBalikCreateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	if request.IdAdminPuskesmas > 0 {
		if err := s.PasienRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, &entity.Pasien{}, request.IdPasien, request.IdAdminPuskesmas, constant.StatusPasienAktif); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.PasienRepository.FindByIdAndStatus(tx, &entity.Pasien{}, request.IdPasien, constant.StatusPasienAktif); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	}

	kontrolBalik := new(entity.KontrolBalik)
	kontrolBalik.IdPasien = request.IdPasien
	kontrolBalik.TanggalKontrol = request.TanggalKontrol
	kontrolBalik.Status = constant.StatusKontrolBalikMenunggu

	if err := s.KontrolBalikRepository.Create(tx, kontrolBalik); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *KontrolBalikService) Update(ctx context.Context, request *model.KontrolBalikUpdateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	kontrolBalik := new(entity.KontrolBalik)
	if request.IdAdminPuskesmas > 0 {
		if err := s.KontrolBalikRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, kontrolBalik, request.ID, request.IdAdminPuskesmas, constant.StatusKontrolBalikMenunggu); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.KontrolBalikRepository.FindByIdAndStatus(tx, kontrolBalik, request.ID, constant.StatusKontrolBalikMenunggu); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	}

	if request.IdAdminPuskesmas > 0 {
		if err := s.PasienRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, &entity.Pasien{}, request.IdPasien, request.IdAdminPuskesmas, constant.StatusPasienAktif); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.PasienRepository.FindByIdAndStatus(tx, &entity.Pasien{}, request.IdPasien, constant.StatusPasienAktif); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	}

	kontrolBalik.IdPasien = request.IdPasien
	kontrolBalik.TanggalKontrol = request.TanggalKontrol

	if err := s.KontrolBalikRepository.Update(tx, kontrolBalik); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *KontrolBalikService) Delete(ctx context.Context, request *model.KontrolBalikDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	kontrolBalik := new(entity.KontrolBalik)
	if request.IdAdminPuskesmas > 0 {
		if err := s.KontrolBalikRepository.FindByIdAndIdAdminPuskesmasAndStatusOrStatus(tx, kontrolBalik, request.ID, request.IdAdminPuskesmas, constant.StatusKontrolBalikBatal, constant.StatusKontrolBalikSelesai); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.KontrolBalikRepository.FindByIdAndStatusOrStatus(tx, kontrolBalik, request.ID, constant.StatusKontrolBalikBatal, constant.StatusKontrolBalikSelesai); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	}

	if err := s.KontrolBalikRepository.Delete(tx, kontrolBalik); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *KontrolBalikService) Batal(ctx context.Context, request *model.KontrolBalikBatalRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	kontrolBalik := new(entity.KontrolBalik)
	if request.IdAdminPuskesmas > 0 {
		if err := s.KontrolBalikRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, kontrolBalik, request.ID, request.IdAdminPuskesmas, constant.StatusKontrolBalikMenunggu); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.KontrolBalikRepository.FindByIdAndStatus(tx, kontrolBalik, request.ID, constant.StatusKontrolBalikMenunggu); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	}

	kontrolBalik.Status = constant.StatusKontrolBalikBatal

	if err := s.KontrolBalikRepository.Update(tx, kontrolBalik); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *KontrolBalikService) Selesai(ctx context.Context, request *model.KontrolBalikSelesaiRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	kontrolBalik := new(entity.KontrolBalik)
	if request.IdAdminPuskesmas > 0 {
		if err := s.KontrolBalikRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, kontrolBalik, request.ID, request.IdAdminPuskesmas, constant.StatusKontrolBalikMenunggu); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.KontrolBalikRepository.FindByIdAndStatus(tx, kontrolBalik, request.ID, constant.StatusKontrolBalikMenunggu); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	}

	kontrolBalik.Status = constant.StatusKontrolBalikSelesai

	if err := s.KontrolBalikRepository.Update(tx, kontrolBalik); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}
