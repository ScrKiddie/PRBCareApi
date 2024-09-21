package service

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log/slog"
	"prb_care_api/internal/adapter"
	"prb_care_api/internal/entity"
	"prb_care_api/internal/model"
	"prb_care_api/internal/repository"
	"strings"
	"sync"
	"time"
)

type ArtikelService struct {
	DB                       *gorm.DB
	ArtikelRepository        *repository.ArtikelRepository
	AdminPuskesmasRepository *repository.AdminPuskesmasRepository
	FileRepository           *repository.FileRepository
	FileAdapter              *adapter.FileAdapter
	Validator                *validator.Validate
	Config                   *viper.Viper
}

func NewArtikelService(
	db *gorm.DB,
	artikelRepository *repository.ArtikelRepository,
	adminPuskesmasRepository *repository.AdminPuskesmasRepository,
	fileRepository *repository.FileRepository,
	fileAdapter *adapter.FileAdapter,
	validator *validator.Validate,
	config *viper.Viper,
) *ArtikelService {
	return &ArtikelService{
		DB:                       db,
		ArtikelRepository:        artikelRepository,
		AdminPuskesmasRepository: adminPuskesmasRepository,
		FileRepository:           fileRepository,
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

func (s *ArtikelService) Create(ctx context.Context, request *model.ArtikelCreateRequest, file *model.File) error {
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

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(request.Isi))
	if err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	var newFileNames []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	doc.Find("img").Each(func(i int, g *goquery.Selection) {
		src, exists := g.Attr("src")
		if exists && strings.HasPrefix(src, "data:image/") {
			wg.Add(1)
			go func(src string, g *goquery.Selection) {
				defer wg.Done()
				imgName, saveErr := s.FileAdapter.StoreImageFromBase64(s.Config.GetString("dir.pict"), src)
				mu.Lock()
				defer mu.Unlock()
				if saveErr != nil {
					slog.Error(saveErr.Error())
					newFileNames = append(newFileNames, "")
					g.SetAttr("src", "")
					return
				}
				newFileNames = append(newFileNames, imgName.Name)
				g.SetAttr("src", imgName.Name)
			}(src, g)
		}
	})

	wg.Wait()

	updatedContent, err := doc.Html()
	if err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	artikel := &entity.Artikel{
		Judul:            request.Judul,
		Ringkasan:        request.Ringkasan,
		Isi:              updatedContent,
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

	for _, imgName := range newFileNames {
		if imgName == "" {
			continue
		}
		fileRecord := &entity.File{
			IdArtikel: artikel.ID,
			File:      imgName,
		}
		if err := s.FileRepository.Create(tx, fileRecord); err != nil {
			slog.Error(err.Error())
			return fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}
	return nil
}

func (s *ArtikelService) Update(ctx context.Context, request *model.ArtikelUpdateRequest, file *model.File) error {
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

	// parsing dan handling konten artikel
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(request.Isi))
	if err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	var newFileNames []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	doc.Find("img").Each(func(i int, g *goquery.Selection) {
		src, exists := g.Attr("src")
		if exists {
			wg.Add(1)
			go func(src string, g *goquery.Selection) {
				defer wg.Done()
				if strings.HasPrefix(src, "data:image/") {
					imgName, saveErr := s.FileAdapter.StoreImageFromBase64(s.Config.GetString("dir.pict"), src)
					mu.Lock()
					defer mu.Unlock()
					if saveErr != nil {
						slog.Error(saveErr.Error())
						newFileNames = append(newFileNames, "")
						g.SetAttr("src", "")
						return
					}
					newFileNames = append(newFileNames, imgName.Name)
					g.SetAttr("src", imgName.Name)
				} else {
					mu.Lock()
					defer mu.Unlock()
					newFileNames = append(newFileNames, src)
				}
			}(src, g)
		}
	})
	wg.Wait()

	updatedContent, err := doc.Html()
	if err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	artikel.Judul = request.Judul
	artikel.Ringkasan = request.Ringkasan
	artikel.Isi = updatedContent
	artikel.IdAdminPuskesmas = request.IdAdminPuskesmas

	if storedFile != nil {
		if artikel.Banner != "" {
			deletedFile := new(model.File)
			deletedFile.Name = artikel.Banner
			s.FileAdapter.DeleteFileAsync(s.Config.GetString("dir.pict"), deletedFile)
		}
		artikel.Banner = storedFile.Name
	}

	// simpan update artikel ke database
	if err := s.ArtikelRepository.Update(tx, artikel); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}
	// ambil file-file yang ada di database terkait artikel ini
	var storedFiles []entity.File
	if err := s.FileRepository.SearchByIdArtikel(tx, &storedFiles, artikel.ID); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	// cek perbedaan file baru dengan file yang ada di database
	existingFileMap := make(map[string]bool)
	for _, f := range storedFiles {
		existingFileMap[f.File] = true
	}

	for _, imgName := range newFileNames {
		if imgName == "" {
			continue
		}

		// jika gambar belum ada di database, tambahkan
		if !existingFileMap[imgName] {
			fileRecord := &entity.File{
				IdArtikel: artikel.ID,
				File:      imgName,
			}
			if err := s.FileRepository.Create(tx, fileRecord); err != nil {
				slog.Error(err.Error())
				return fiber.ErrInternalServerError
			}
		}
		// hapus dari map karena masih digunakan
		delete(existingFileMap, imgName)
	}

	// hapus file yang tidak ada dalam daftar baru
	for imgName := range existingFileMap {
		// cari file yang cocok dengan imgName dari storedFiles untuk mendapatkan ID
		var fileToDelete *entity.File
		for _, storedFile := range storedFiles {
			if storedFile.File == imgName {
				fileToDelete = &storedFile
				break
			}
		}

		if fileToDelete != nil {
			// hapus file dari database menggunakan ID dan nama file
			if err := s.FileRepository.Delete(tx, fileToDelete); err != nil {
				slog.Error(err.Error())
				return fiber.ErrInternalServerError
			}

			// hapus file dari storage
			fileModel := new(model.File)
			fileModel.Name = imgName
			s.FileAdapter.DeleteFileAsync(s.Config.GetString("dir.pict"), fileModel)
		}
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

	var files []entity.File
	if err := s.FileRepository.SearchByIdArtikel(tx, &files, artikel.ID); err != nil {
		slog.Error(err.Error())
		return fiber.ErrInternalServerError
	}

	for _, file := range files {
		if err := s.FileRepository.Delete(tx, &file); err != nil {
			slog.Error(err.Error())
			return fiber.ErrInternalServerError
		}

		fileModel := &model.File{Name: file.File}
		s.FileAdapter.DeleteFileAsync(s.Config.GetString("dir.pict"), fileModel)
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
