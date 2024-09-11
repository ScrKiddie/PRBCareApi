package repository

import (
	"gorm.io/gorm"
	"prb_care_api/internal/entity"
)

type PasienRepository struct {
	Repository[entity.Pasien]
}

func NewPasienRepository() *PasienRepository {
	return &PasienRepository{}
}

func (r *PasienRepository) Search(db *gorm.DB, pasien *[]entity.Pasien, status string) error {
	query := db.Preload("AdminPuskesmas").Preload("Pengguna")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	return query.Find(pasien).Error
}
func (r *PasienRepository) SearchAsAdminPuskesmas(db *gorm.DB, pasien *[]entity.Pasien, idAdminPuskesmas int32, status string) error {
	query := db.Where("id_admin_puskesmas = ?", idAdminPuskesmas).Preload("Pengguna")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	return query.Find(pasien).Error
}
func (r *PasienRepository) SearchAsPengguna(db *gorm.DB, pasien *[]entity.Pasien, idPengguna int32, status string) error {
	query := db.Where("id_pengguna = ?", idPengguna).Preload("AdminPuskesmas")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	return query.Find(pasien).Error
}
func (r *PasienRepository) FindByIdAndStatus(db *gorm.DB, pasien *entity.Pasien, id int32, status string) error {
	return db.Where("id = ?", id).Where("status = ?", status).First(pasien).Error
}
func (r *PasienRepository) FindByIdAndIdAdminPuskesmasAndStatus(db *gorm.DB, pasien *entity.Pasien, id int32, idAdminPuskesmas int32, status string) error {
	return db.Where("id = ?", id).Where("id_admin_puskesmas = ?", idAdminPuskesmas).Where("status = ?", status).First(pasien).Error
}
func (r *PasienRepository) FindByIdAdminPuskesmas(db *gorm.DB, pasien *entity.Pasien, idAdminPuskesmas int32) error {
	return db.Where("id_admin_puskesmas = ?", idAdminPuskesmas).First(pasien).Error
}
func (r *PasienRepository) FindByIdPengguna(db *gorm.DB, pasien *entity.Pasien, idPengguna int32) error {
	return db.Where("id_pengguna = ?", idPengguna).First(pasien).Error
}
