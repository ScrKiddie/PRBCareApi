package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/spf13/viper"
	"prb_care_api/internal/controller"
)

type Config struct {
	App                       *fiber.App
	AuthMiddleware            fiber.Handler
	AdminSuperController      *controller.AdminSuperController
	AdminPuskesmasController  *controller.AdminPuskesmasController
	AdminApotekController     *controller.AdminApotekController
	PenggunaController        *controller.PenggunaController
	ObatController            *controller.ObatController
	PasienController          *controller.PasienController
	KontrolBalikController    *controller.KontrolBalikController
	PengambilanObatController *controller.PengambilanObatController
	ArtikelController         *controller.ArtikelController
	Config                    *viper.Viper
}

func (c *Config) SetupGuestRoute() {
	c.App.Use(cors.New(cors.Config{
		AllowOrigins: c.Config.GetStringSlice("web.cors.origins"),
		AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	c.App.Post("/api/admin-super/login", c.AdminSuperController.Login)
	c.App.Post("/api/admin-puskesmas/login", c.AdminPuskesmasController.Login)
	c.App.Post("/api/admin-apotek/login", c.AdminApotekController.Login)
	c.App.Post("/api/pengguna/login", c.PenggunaController.Login)
	c.App.Post("/api/pengguna/register", c.PenggunaController.Register)

}

func (c *Config) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)

	c.App.Patch("/api/admin-super/current/password", c.AdminSuperController.PasswordUpdate)

	c.App.Get("/api/admin-puskesmas", c.AdminPuskesmasController.List)
	c.App.Get("/api/admin-puskesmas/current", c.AdminPuskesmasController.Current)
	c.App.Patch("/api/admin-puskesmas/current", c.AdminPuskesmasController.CurrentProfileUpdate)
	c.App.Patch("/api/admin-puskesmas/current/password", c.AdminPuskesmasController.CurrentPasswordUpdate)
	c.App.Get("/api/admin-puskesmas/:id", c.AdminPuskesmasController.Get)
	c.App.Post("/api/admin-puskesmas", c.AdminPuskesmasController.Create)
	c.App.Patch("/api/admin-puskesmas/:id", c.AdminPuskesmasController.Update)
	c.App.Delete("/api/admin-puskesmas/:id", c.AdminPuskesmasController.Delete)

	c.App.Get("/api/admin-apotek", c.AdminApotekController.List)
	c.App.Get("/api/admin-apotek/current", c.AdminApotekController.Current)
	c.App.Patch("/api/admin-apotek/current", c.AdminApotekController.CurrentProfileUpdate)
	c.App.Patch("/api/admin-apotek/current/password", c.AdminApotekController.CurrentPasswordUpdate)
	c.App.Get("/api/admin-apotek/:id", c.AdminApotekController.Get)
	c.App.Post("/api/admin-apotek", c.AdminApotekController.Create)
	c.App.Patch("/api/admin-apotek/:id", c.AdminApotekController.Update)
	c.App.Delete("/api/admin-apotek/:id", c.AdminApotekController.Delete)

	c.App.Get("/api/pengguna/current", c.PenggunaController.Current)
	c.App.Patch("/api/pengguna/current", c.PenggunaController.CurrentProfileUpdate)
	c.App.Patch("/api/pengguna/current/password", c.PenggunaController.CurrentPasswordUpdate)
	c.App.Patch("/api/pengguna/current/perangkat", c.PenggunaController.CurrentTokenPerangkatUpdate)
	c.App.Get("/api/pengguna", c.PenggunaController.List)
	c.App.Get("/api/pengguna/:id", c.PenggunaController.Get)
	c.App.Post("/api/pengguna", c.PenggunaController.Create)
	c.App.Patch("/api/pengguna/:id", c.PenggunaController.Update)
	c.App.Delete("/api/pengguna/:id", c.PenggunaController.Delete)

	c.App.Get("/api/obat", c.ObatController.List)
	c.App.Get("/api/obat/:id", c.ObatController.Get)
	c.App.Post("/api/obat", c.ObatController.Create)
	c.App.Patch("/api/obat/:id", c.ObatController.Update)
	c.App.Delete("/api/obat/:id", c.ObatController.Delete)

	c.App.Get("/api/pasien", c.PasienController.Search)
	c.App.Get("/api/pasien/:id", c.PasienController.Get)
	c.App.Post("/api/pasien", c.PasienController.Create)
	c.App.Patch("/api/pasien/:id", c.PasienController.Update)
	c.App.Delete("/api/pasien/:id", c.PasienController.Delete)
	c.App.Patch("/api/pasien/:id/selesai", c.PasienController.Selesai)

	c.App.Get("/api/kontrol-balik", c.KontrolBalikController.Search)
	c.App.Get("/api/kontrol-balik/:id", c.KontrolBalikController.Get)
	c.App.Post("/api/kontrol-balik", c.KontrolBalikController.Create)
	c.App.Patch("/api/kontrol-balik/:id", c.KontrolBalikController.Update)
	c.App.Delete("/api/kontrol-balik/:id", c.KontrolBalikController.Delete)
	c.App.Patch("/api/kontrol-balik/:id/selesai", c.KontrolBalikController.Selesai)
	c.App.Patch("/api/kontrol-balik/:id/batal", c.KontrolBalikController.Batal)

	c.App.Get("/api/pengambilan-obat", c.PengambilanObatController.Search)
	c.App.Get("/api/pengambilan-obat/:id", c.PengambilanObatController.Get)
	c.App.Post("/api/pengambilan-obat", c.PengambilanObatController.Create)
	c.App.Patch("/api/pengambilan-obat/:id", c.PengambilanObatController.Update)
	c.App.Delete("/api/pengambilan-obat/:id", c.PengambilanObatController.Delete)
	c.App.Patch("/api/pengambilan-obat/:id/batal", c.PengambilanObatController.Batal)
	c.App.Patch("/api/pengambilan-obat/:id/diambil", c.PengambilanObatController.Diambil)

	c.App.Get("/api/artikel", c.ArtikelController.Search)
	c.App.Get("/api/artikel/:id", c.ArtikelController.Get)
	c.App.Post("/api/artikel", c.ArtikelController.Create)
	c.App.Patch("/api/artikel/:id", c.ArtikelController.Update)
	c.App.Delete("/api/artikel/:id", c.ArtikelController.Delete)
}

func (c *Config) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}
