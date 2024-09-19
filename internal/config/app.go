package config

import (
	"github.com/go-playground/mold/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/client"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"prb_care_api/internal/adapter"
	"prb_care_api/internal/controller"
	"prb_care_api/internal/middleware"
	"prb_care_api/internal/repository"
	"prb_care_api/internal/route"
	"prb_care_api/internal/service"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Client   *client.Client
	Validate *validator.Validate
	Config   *viper.Viper
	Modifier *mold.Transformer
}

func Bootstrap(config *BootstrapConfig) {

	adminSuperRepository := repository.NewAdminSuperRepository()
	adminPuskesmasRepository := repository.NewAdminPuskesmasRepository()
	adminApotekRepository := repository.NewAdminApotekRepository()
	penggunaRepository := repository.NewPenggunaRepository()
	obatRepository := repository.NewObatRepository()
	pasienRepository := repository.NewPasienRepository()
	kontrolBalikRepository := repository.NewKontrolBalikRepository()
	pengambilanObatRepository := repository.NewPengambilanObatRepository()
	artikelRepository := repository.NewArtikelRepository()

	captchaAdapter := adapter.NewCaptcha(config.Client)
	fileAdapter := adapter.NewFileAdapter()

	adminSuperService := service.NewAdminSuperService(config.DB, adminSuperRepository, captchaAdapter, config.Validate, config.Config)
	adminPuskesmasService := service.NewAdminPuskesmasService(config.DB, adminPuskesmasRepository, pasienRepository, captchaAdapter, config.Validate, config.Config)
	adminApotekService := service.NewAdminApotekService(config.DB, adminApotekRepository, obatRepository, config.Validate, captchaAdapter, config.Config)
	penggunaService := service.NewPenggunaService(config.DB, penggunaRepository, pasienRepository, config.Validate, captchaAdapter, config.Config)
	obatService := service.NewObatService(config.DB, obatRepository, adminApotekRepository, pengambilanObatRepository, config.Validate)
	pasienService := service.NewPasienService(config.DB, pasienRepository, adminPuskesmasRepository, penggunaRepository, kontrolBalikRepository, pengambilanObatRepository, config.Validate)
	kontrolBalikService := service.NewKontrolBalikService(config.DB, kontrolBalikRepository, pasienRepository, config.Validate)
	pengambilanObatService := service.NewPengambilanObatService(config.DB, pengambilanObatRepository, pasienRepository, obatRepository, config.Validate)
	artikelSevice := service.NewArtikelService(config.DB, artikelRepository, adminPuskesmasRepository, fileAdapter, config.Validate, config.Config)

	adminSuperController := controller.NewAdminSuperController(adminSuperService)
	adminPuskesmasController := controller.NewAdminPuskesmasController(adminPuskesmasService, config.Modifier)
	adminApotekController := controller.NewAdminApotekController(adminApotekService, config.Modifier)
	penggunaController := controller.NewPenggunaController(penggunaService, config.Modifier)
	obatController := controller.NewObatController(obatService, config.Modifier)
	pasienController := controller.NewPasienController(pasienService, config.Modifier)
	kontrolBalikController := controller.NewKontrolBalikController(kontrolBalikService, config.Modifier)
	pengambilanObatController := controller.NewPengambilanObatController(pengambilanObatService)
	artikelController := controller.NewArtikelController(artikelSevice, config.Modifier)

	authMiddleware := middleware.AuthMiddleware(config.Config, adminSuperService, adminPuskesmasService, adminApotekService, penggunaService)

	route := route.Config{
		App:                       config.App,
		AuthMiddleware:            authMiddleware,
		AdminSuperController:      adminSuperController,
		AdminPuskesmasController:  adminPuskesmasController,
		AdminApotekController:     adminApotekController,
		PenggunaController:        penggunaController,
		ObatController:            obatController,
		PasienController:          pasienController,
		KontrolBalikController:    kontrolBalikController,
		PengambilanObatController: pengambilanObatController,
		ArtikelController:         artikelController,
		Config:                    config.Config,
	}
	route.Setup()

}
