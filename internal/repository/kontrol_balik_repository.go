package repository

import (
	"gorm.io/gorm"
	"prb_care_api/internal/entity"
)

type KontrolBalikRepository struct {
	Repository[entity.KontrolBalik]
}

func NewKontrolBalikRepository() *KontrolBalikRepository {
	return &KontrolBalikRepository{}
}

func (r *KontrolBalikRepository) Search(db *gorm.DB, kontrolBalik *[]entity.KontrolBalik, status string) error {
	query := db
	if status != "" {
		query = query.Where("status = ?", status)
	}
	return query.Preload("Pasien.AdminPuskesmas").Preload("Pasien.Pengguna").Find(&kontrolBalik).Error
}
func (r *KontrolBalikRepository) SearchAsAdminPuskesmas(db *gorm.DB, kontrolBalik *[]entity.KontrolBalik, idAdminPuskesmas int, status string) error {
	query := db.Joins("JOIN pasien ON pasien.id = kontrol_balik.id_pasien").
		Where("pasien.id_admin_puskesmas = ?", idAdminPuskesmas)
	if status != "" {
		query = query.Where("kontrol_balik.status = ?", status)
	}
	return query.Preload("Pasien.AdminPuskesmas").Preload("Pasien.Pengguna").Find(&kontrolBalik).Error
}
func (r *KontrolBalikRepository) SearchAsPengguna(db *gorm.DB, kontrolBalik *[]entity.KontrolBalik, idPengguna int, status string) error {
	query := db.Joins("JOIN pasien ON pasien.id = kontrol_balik.id_pasien").
		Where("pasien.id_pengguna = ?", idPengguna)
	if status != "" {
		query = query.Where("kontrol_balik.status = ?", status)
	}
	return query.Preload("Pasien.AdminPuskesmas").Preload("Pasien.Pengguna").Find(&kontrolBalik).Error
}
func (r *KontrolBalikRepository) FindByIdAndStatus(db *gorm.DB, kontrolBalik *entity.KontrolBalik, id int, status string) error {
	return db.Where("id = ?", id).
		Where("status = ?", status).
		First(&kontrolBalik).Error
}
func (r *KontrolBalikRepository) FindByIdAndIdAdminPuskesmasAndStatus(db *gorm.DB, kontrolBalik *entity.KontrolBalik, id int, idAdminPuskesmas int, status string) error {
	return db.Joins("JOIN pasien ON pasien.id = kontrol_balik.id_pasien").
		Where("kontrol_balik.id = ?", id).
		Where("pasien.id_admin_puskesmas = ?", idAdminPuskesmas).
		Where("kontrol_balik.status = ?", status).
		First(&kontrolBalik).Error
}
func (r *KontrolBalikRepository) FindByIdAndIdAdminPuskesmasAndStatusOrStatus(db *gorm.DB, kontrolBalik *entity.KontrolBalik, id int, idAdminPuskesmas int, status1 string, status2 string) error {
	return db.Joins("JOIN pasien ON pasien.id = kontrol_balik.id_pasien").
		Where("kontrol_balik.id = ?", id).
		Where("pasien.id_admin_puskesmas = ?", idAdminPuskesmas).
		Where("kontrol_balik.status = ? OR kontrol_balik.status = ?", status1, status2).
		First(&kontrolBalik).Error
}
func (r *KontrolBalikRepository) FindByIdAndStatusOrStatus(db *gorm.DB, kontrolBalik *entity.KontrolBalik, id int, status1 string, status2 string) error {
	return db.
		Where("id = ?", id).
		Where("status = ? OR status = ?", status1, status2).
		First(&kontrolBalik).Error
}
func (r *KontrolBalikRepository) FindByIdPasienAndStatus(db *gorm.DB, kontrolBalik *entity.KontrolBalik, idPasien int, status string) error {
	return db.Where("id_pasien = ?", idPasien).
		Where("status = ?", status).
		First(&kontrolBalik).Error
}
func (r *KontrolBalikRepository) FindByIdPasien(db *gorm.DB, kontrolBalik *entity.KontrolBalik, idPasien int) error {
	return db.Where("id_pasien = ?", idPasien).
		First(&kontrolBalik).Error
}
