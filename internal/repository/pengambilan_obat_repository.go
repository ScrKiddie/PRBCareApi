package repository

import (
	"gorm.io/gorm"
	"prb_care_api/internal/entity"
)

type PengambilanObatRepository struct {
	Repository[entity.PengambilanObat]
}

func NewPengambilanObatRepository() *PengambilanObatRepository {
	return &PengambilanObatRepository{}
}

func (r *PengambilanObatRepository) Search(db *gorm.DB, pengambilanObat *[]entity.PengambilanObat, status string) error {
	query := db
	if status != "" {
		query = query.Where("status = ?", status)
	}
	return query.Preload("Pasien.AdminPuskesmas").
		Preload("Pasien.Pengguna").
		Preload("Obat.AdminApotek").
		Find(&pengambilanObat).Error
}
func (r *PengambilanObatRepository) SearchAsAdminPuskesmas(db *gorm.DB, pengambilanObat *[]entity.PengambilanObat, idAdminPuskesmas int32, status string) error {
	query := db.Joins("JOIN pasien ON pasien.id = pengambilan_obat.id_pasien").
		Where("pasien.id_admin_puskesmas = ?", idAdminPuskesmas)
	if status != "" {
		query = query.Where("pengambilan_obat.status = ?", status)
	}
	return query.
		Preload("Pasien.Pengguna").
		Preload("Obat.AdminApotek").
		Find(&pengambilanObat).Error
}
func (r *PengambilanObatRepository) SearchAsAdminApotek(db *gorm.DB, pengambilanObat *[]entity.PengambilanObat, idAdminApotek int32, status string) error {
	query := db.Joins("JOIN obat ON obat.id = pengambilan_obat.id_obat").
		Where("obat.id_admin_apotek = ?", idAdminApotek)
	if status != "" {
		query = query.Where("pengambilan_obat.status = ?", status)
	}
	return query.Preload("Pasien.AdminPuskesmas").
		Preload("Pasien.Pengguna").
		Find(&pengambilanObat).Error
}
func (r *PengambilanObatRepository) SearchAsPengguna(db *gorm.DB, pengambilanObat *[]entity.PengambilanObat, idPengguna int32, status string) error {
	query := db.Joins("JOIN pasien ON pasien.id = pengambilan_obat.id_pasien").
		Where("pasien.id_pengguna = ?", idPengguna)
	if status != "" {
		query = query.Where("pengambilan_obat.status = ?", status)
	}
	return query.Preload("Pasien.AdminPuskesmas").
		Preload("Obat.AdminApotek").
		Find(&pengambilanObat).Error
}
func (r *PengambilanObatRepository) FindByIdAndStatus(db *gorm.DB, pengambilanObat *entity.PengambilanObat, id int32, status string) error {
	return db.Where("id = ?", id).
		Where("status = ?", status).
		First(&pengambilanObat).Error
}
func (r *PengambilanObatRepository) FindByIdAndIdAdminPuskesmasAndStatus(db *gorm.DB, pengambilanObat *entity.PengambilanObat, id int32, idAdminPuskesmas int32, status string) error {
	return db.Joins("JOIN pasien ON pasien.id = pengambilan_obat.id_pasien").
		Where("pengambilan_obat.id = ?", id).
		Where("pasien.id_admin_puskesmas = ?", idAdminPuskesmas).
		Where("pengambilan_obat.status = ?", status).
		First(&pengambilanObat).Error
}
func (r *PengambilanObatRepository) FindByIdAndIdAdminApotekAndStatus(db *gorm.DB, pengambilanObat *entity.PengambilanObat, id int32, idAdminApotek int32, status string) error {
	return db.Joins("JOIN obat ON obat.id = pengambilan_obat.id_obat").
		Where("pengambilan_obat.id = ?", id).
		Where("obat.id_admin_apotek = ?", idAdminApotek).
		Where("pengambilan_obat.status = ?", status).
		First(&pengambilanObat).Error
}
func (r *PengambilanObatRepository) FindByIdAndStatusOrStatus(db *gorm.DB, pengambilanObat *entity.PengambilanObat, id int32, status1 string, status2 string) error {
	return db.Where("id = ?", id).
		Where("status = ? OR status = ?", status1, status2).
		First(&pengambilanObat).Error
}
func (r *PengambilanObatRepository) FindByIdAndIdAdminPuskesmasAndStatusOrStatus(db *gorm.DB, pengambilanObat *entity.PengambilanObat, id int32, idAdminPuskesmas int32, status1 string, status2 string) error {
	return db.Joins("JOIN pasien ON pasien.id = pengambilan_obat.id_pasien").
		Where("pengambilan_obat.id = ?", id).
		Where("pasien.id_admin_puskesmas = ?", idAdminPuskesmas).
		Where("pengambilan_obat.status = ? OR pengambilan_obat.status = ?", status1, status2).
		First(&pengambilanObat).Error
}
func (r *PengambilanObatRepository) FindByIdPasienAndStatus(db *gorm.DB, pengambilanObat *entity.PengambilanObat, idPasien int32, status string) error {
	return db.
		Where("id_pasien = ?", idPasien).
		Where("status = ?", status).
		First(&pengambilanObat).Error
}
func (r *PengambilanObatRepository) FindByIdObat(db *gorm.DB, pengambilanObat *entity.PengambilanObat, idObat int32) error {
	return db.Where("id_obat = ?", idObat).First(&pengambilanObat).Error
}
func (r *PengambilanObatRepository) FindByIdPasien(db *gorm.DB, pengambilanObat *entity.PengambilanObat, idPasien int32) error {
	return db.Where("id_pasien = ?", idPasien).First(&pengambilanObat).Error
}
