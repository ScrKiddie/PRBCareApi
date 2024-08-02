package service

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"log"
	"prb_care_api/internal/constant"
	"prb_care_api/internal/entity"
	"prb_care_api/internal/model"
	"prb_care_api/internal/repository"
)

type PengambilanObatService struct {
	DB                        *gorm.DB
	PengambilanObatRepository *repository.PengambilanObatRepository
	PasienRepository          *repository.PasienRepository
	ObatRepository            *repository.ObatRepository
	Validator                 *validator.Validate
}

func NewPengambilanObatService(
	db *gorm.DB,
	pengambilanObatRepository *repository.PengambilanObatRepository,
	pasienRepository *repository.PasienRepository,
	obatRepository *repository.ObatRepository,
	validator *validator.Validate,
) *PengambilanObatService {
	return &PengambilanObatService{db, pengambilanObatRepository, pasienRepository, obatRepository, validator}
}

func (s *PengambilanObatService) Search(ctx context.Context, request *model.PengambilanObatSearchRequest) (*[]model.PengambilanObatResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	pengambilanObat := new([]entity.PengambilanObat)
	if request.IdPengguna > 0 {
		if err := s.PengambilanObatRepository.SearchAsPengguna(tx, pengambilanObat, request.IdPengguna, request.Status); err != nil {
			log.Println(err.Error())
			return nil, fiber.ErrInternalServerError
		}
	} else if request.IdAdminPuskesmas > 0 {
		if err := s.PengambilanObatRepository.SearchAsAdminPuskesmas(tx, pengambilanObat, request.IdAdminPuskesmas, request.Status); err != nil {
			log.Println(err.Error())
			return nil, fiber.ErrInternalServerError
		}
	} else if request.IdAdminApotek > 0 {
		if err := s.PengambilanObatRepository.SearchAsAdminApotek(tx, pengambilanObat, request.IdAdminApotek, request.Status); err != nil {
			log.Println(err.Error())
			return nil, fiber.ErrInternalServerError
		}
	} else {
		if err := s.PengambilanObatRepository.Search(tx, pengambilanObat, request.Status); err != nil {
			log.Println(err.Error())
			return nil, fiber.ErrInternalServerError
		}
	}

	var response []model.PengambilanObatResponse
	for _, p := range *pengambilanObat {
		response = append(response, model.PengambilanObatResponse{
			ID:   p.ID,
			Resi: p.Resi,
			PasienResponse: &model.PasienResponse{
				ID:           p.Pasien.ID,
				NoRekamMedis: p.Pasien.NoRekamMedis,
				Pengguna: &model.PenggunaResponse{
					ID:              p.Pasien.Pengguna.ID,
					NamaLengkap:     p.Pasien.Pengguna.NamaLengkap,
					Telepon:         p.Pasien.Pengguna.Telepon,
					TeleponKeluarga: p.Pasien.Pengguna.TeleponKeluarga,
					Alamat:          p.Pasien.Pengguna.Alamat,
				},
				AdminPuskesmas: &model.AdminPuskesmasResponse{
					ID:            p.Pasien.AdminPuskesmas.ID,
					NamaPuskesmas: p.Pasien.AdminPuskesmas.NamaPuskesmas,
					Telepon:       p.Pasien.AdminPuskesmas.Telepon,
					Alamat:        p.Pasien.AdminPuskesmas.Alamat,
				},
				BeratBadan:     p.Pasien.BeratBadan,
				TinggiBadan:    p.Pasien.TinggiBadan,
				TekananDarah:   p.Pasien.TekananDarah,
				DenyutNadi:     p.Pasien.DenyutNadi,
				HasilLab:       p.Pasien.HasilLab,
				HasilEkg:       p.Pasien.HasilEkg,
				TanggalPeriksa: p.Pasien.TanggalPeriksa,
				Status:         p.Pasien.Status,
			},
			Obat: &model.ObatResponse{
				ID:       p.Obat.ID,
				NamaObat: p.Obat.NamaObat,
				Jumlah:   p.Obat.Jumlah,
				AdminApotek: &model.AdminApotekResponse{
					ID:         p.Obat.AdminApotek.ID,
					NamaApotek: p.Obat.AdminApotek.NamaApotek,
					Telepon:    p.Obat.AdminApotek.Telepon,
					Alamat:     p.Obat.AdminApotek.Alamat,
				},
			},
			Jumlah:             p.Jumlah,
			TanggalPengambilan: p.TanggalPengambilan,
			Status:             p.Status,
		})
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return &response, nil
}

func (s *PengambilanObatService) Get(ctx context.Context, request *model.PengambilanObatGetRequest) (*model.PengambilanObatResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrBadRequest
	}

	pengambilanObat := new(entity.PengambilanObat)
	if request.IdAdminPuskesmas > 0 {
		if err := s.PengambilanObatRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, pengambilanObat, request.ID, request.IdAdminPuskesmas, constant.StatusPengambilanObatMenunggu); err != nil {
			log.Println(err.Error())
			return nil, fiber.ErrNotFound
		}
	} else if err := s.PengambilanObatRepository.FindByIdAndStatus(tx, pengambilanObat, request.ID, constant.StatusPengambilanObatMenunggu); err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return nil, fiber.ErrInternalServerError
	}

	response := new(model.PengambilanObatResponse)
	response.ID = pengambilanObat.ID
	response.IdObat = pengambilanObat.IdObat
	response.IdPasien = pengambilanObat.IdPasien
	response.Jumlah = pengambilanObat.Jumlah
	response.TanggalPengambilan = pengambilanObat.TanggalPengambilan
	return response, nil
}

func (s *PengambilanObatService) Create(ctx context.Context, request *model.PengambilanObatCreateRequest) error {
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

	obat := new(entity.Obat)
	if err := s.ObatRepository.FindByIdAndLockForUpdate(tx, obat, request.IdObat); err != nil {
		log.Println(err.Error())
		return fiber.ErrNotFound
	}
	obat.Jumlah -= request.Jumlah
	if obat.Jumlah < 0 {
		return fiber.NewError(fiber.StatusConflict, "Jumlah obat melebihi persediaan apotek")
	}

	pengambilanObat := new(entity.PengambilanObat)
	pengambilanObat.Resi = fmt.Sprintf("%d%d", request.IdPasien, request.TanggalPengambilan)
	pengambilanObat.IdPasien = request.IdPasien
	pengambilanObat.IdObat = request.IdObat
	pengambilanObat.Jumlah = request.Jumlah
	pengambilanObat.TanggalPengambilan = request.TanggalPengambilan
	pengambilanObat.Status = constant.StatusPengambilanObatMenunggu

	if err := s.ObatRepository.Update(tx, obat); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := s.PengambilanObatRepository.Create(tx, pengambilanObat); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *PengambilanObatService) Update(ctx context.Context, request *model.PengambilanObatUpdateRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	pengambilanObat := new(entity.PengambilanObat)
	if request.IdAdminPuskesmas > 0 {
		if err := s.PengambilanObatRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, pengambilanObat, request.ID, request.IdAdminPuskesmas, constant.StatusPengambilanObatMenunggu); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.PengambilanObatRepository.FindByIdAndStatus(tx, pengambilanObat, request.ID, constant.StatusPengambilanObatMenunggu); err != nil {
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

	obatNew := new(entity.Obat)
	if err := s.ObatRepository.FindByIdAndLockForUpdate(tx, obatNew, request.IdObat); err != nil {
		log.Println(err.Error())
		return fiber.ErrNotFound
	}

	obatOld := new(entity.Obat)
	if pengambilanObat.IdObat == request.IdObat {
		obatNew.Jumlah = (obatNew.Jumlah + pengambilanObat.Jumlah) - request.Jumlah
		if obatNew.Jumlah < 0 {
			return fiber.NewError(fiber.StatusConflict, "Jumlah obat melebihi persediaan apotek")
		}
	} else {
		if err := s.ObatRepository.FindByIdAndLockForUpdate(tx, obatOld, pengambilanObat.IdObat); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
		obatOld.Jumlah += pengambilanObat.Jumlah
		obatNew.Jumlah -= request.Jumlah
		if obatNew.Jumlah < 0 {
			return fiber.NewError(fiber.StatusConflict, "Jumlah obat melebihi persediaan apotek")
		}
		if err := s.ObatRepository.Update(tx, obatOld); err != nil {
			log.Println(err.Error())
			return fiber.ErrInternalServerError
		}
	}

	if err := s.ObatRepository.Update(tx, obatNew); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	pengambilanObat.Resi = fmt.Sprintf("%d%d", request.IdPasien, request.TanggalPengambilan)
	pengambilanObat.IdPasien = request.IdPasien
	pengambilanObat.IdObat = request.IdObat
	pengambilanObat.Jumlah = request.Jumlah
	pengambilanObat.TanggalPengambilan = request.TanggalPengambilan

	if err := s.PengambilanObatRepository.Update(tx, pengambilanObat); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *PengambilanObatService) Delete(ctx context.Context, request *model.PengambilanObatDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	pengambilanObat := new(entity.PengambilanObat)
	if request.IdAdminPuskesmas > 0 {
		if err := s.PengambilanObatRepository.FindByIdAndIdAdminPuskesmasAndStatusOrStatus(tx, pengambilanObat, request.ID, request.IdAdminPuskesmas, constant.StatusPengambilanObatBatal, constant.StatusPengambilanObatDiambil); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.PengambilanObatRepository.FindByIdAndStatusOrStatus(tx, pengambilanObat, request.ID, constant.StatusPengambilanObatBatal, constant.StatusPengambilanObatDiambil); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	}

	if err := s.PengambilanObatRepository.Delete(tx, pengambilanObat); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *PengambilanObatService) Batal(ctx context.Context, request *model.PengambilanObatBatalRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	pengambilanObat := new(entity.PengambilanObat)
	if request.IdAdminPuskesmas > 0 {
		if err := s.PengambilanObatRepository.FindByIdAndIdAdminPuskesmasAndStatus(tx, pengambilanObat, request.ID, request.IdAdminPuskesmas, constant.StatusPengambilanObatMenunggu); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.PengambilanObatRepository.FindByIdAndStatus(tx, pengambilanObat, request.ID, constant.StatusPengambilanObatMenunggu); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	}
	obat := new(entity.Obat)
	if err := s.ObatRepository.FindByIdAndLockForUpdate(tx, obat, pengambilanObat.IdObat); err != nil {
		log.Println(err.Error())
		return fiber.ErrNotFound
	}
	obat.Jumlah += pengambilanObat.Jumlah

	pengambilanObat.Status = constant.StatusPengambilanObatBatal

	if err := s.ObatRepository.Update(tx, obat); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := s.PengambilanObatRepository.Update(tx, pengambilanObat); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *PengambilanObatService) Diambil(ctx context.Context, request *model.PengambilanObatDiambilRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validator.Struct(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	pengambilanObat := new(entity.PengambilanObat)
	if request.IdAdminApotek > 0 {
		if err := s.PengambilanObatRepository.FindByIdAndIdAdminApotekAndStatus(tx, pengambilanObat, request.ID, request.IdAdminApotek, constant.StatusPengambilanObatMenunggu); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	} else {
		if err := s.PengambilanObatRepository.FindByIdAndStatus(tx, pengambilanObat, request.ID, constant.StatusPengambilanObatMenunggu); err != nil {
			log.Println(err.Error())
			return fiber.ErrNotFound
		}
	}

	pengambilanObat.Status = constant.StatusPengambilanObatDiambil

	if err := s.PengambilanObatRepository.Update(tx, pengambilanObat); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}
