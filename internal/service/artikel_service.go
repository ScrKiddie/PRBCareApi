package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log/slog"
	"prb_care_api/internal/adapter"
	"prb_care_api/internal/entity"
	"prb_care_api/internal/model"
	"prb_care_api/internal/repository"
	"time"
)

type ArtikelService struct {
	DB                       *gorm.DB
	ArtikelRepository        *repository.ArtikelRepository
	AdminPuskesmasRepository *repository.AdminPuskesmasRepository
	FileAdapter              *adapter.FileAdapter
	Validator                *validator.Validate
	Config                   *viper.Viper
}

func NewArtikelService(
	db *gorm.DB,
	artikelRepository *repository.ArtikelRepository,
	adminPuskesmasRepository *repository.AdminPuskesmasRepository,
	fileAdapter *adapter.FileAdapter,
	validator *validator.Validate,
	config *viper.Viper,
) *ArtikelService {
	return &ArtikelService{
		DB:                       db,
		ArtikelRepository:        artikelRepository,
		AdminPuskesmasRepository: adminPuskesmasRepository,
		FileAdapter:              fileAdapter,
		Validator:                validator,
		Config:                   config,
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
			Banner:           a.Banner,
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
		Banner:           artikel.Banner,
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

func (s *ArtikelService) Create(ctx context.Context, request *model.ArtikelCreateRequest, file *model.FileUpload) error {
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

	var storedFile *model.File
	var err error

	if file != nil {
		if err := s.Validator.Struct(file); err != nil {
			slog.Error(err.Error())
			return fiber.ErrBadRequest
		}
		if file.FileHeader != nil && file.FileHeader.Filename != "" {
			storedFile, err = s.FileAdapter.StoreImage(s.Config.GetString("dir.pict"), file)
			if err != nil {
				slog.Error(err.Error())
				return fiber.ErrInternalServerError
			}
		}
	}

	artikel := &entity.Artikel{
		Judul:            request.Judul,
		Ringkasan:        request.Ringkasan,
		Isi:              request.Isi,
		TanggalPublikasi: time.Now().Unix(),
		IdAdminPuskesmas: request.IdAdminPuskesmas,
	}

	if storedFile != nil {
		artikel.Banner = storedFile.Name
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

func (s *ArtikelService) Update(ctx context.Context, request *model.ArtikelUpdateRequest, file *model.FileUpload) error {
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

	var storedFile *model.File
	var err error

	if file != nil {
		if err := s.Validator.Struct(file); err != nil {
			slog.Error(err.Error())
			return fiber.ErrBadRequest
		}
		if file.FileHeader != nil && file.FileHeader.Filename != "" {
			storedFile, err = s.FileAdapter.StoreImage(s.Config.GetString("dir.pict"), file)
			if err != nil {
				slog.Error(err.Error())
				return fiber.ErrInternalServerError
			}
		}
	}

	artikel.Judul = request.Judul
	artikel.Ringkasan = request.Ringkasan
	artikel.Isi = request.Isi
	artikel.IdAdminPuskesmas = request.IdAdminPuskesmas

	if storedFile != nil {
		if artikel.Banner != "" {
			deletedFile := new(model.File)
			deletedFile.Name = artikel.Banner
			s.FileAdapter.DeleteFileAsync(s.Config.GetString("dir.pict"), deletedFile)
		}
		artikel.Banner = storedFile.Name
	}

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

	if artikel.Banner != "" {
		deletedFile := &model.File{Name: artikel.Banner}
		s.FileAdapter.DeleteFileAsync(s.Config.GetString("dir.pict"), deletedFile)
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
