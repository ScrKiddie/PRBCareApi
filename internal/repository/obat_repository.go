package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"prb_care_api/internal/entity"
)

type ObatRepository struct {
	Repository[entity.Obat]
}

func NewObatRepository() *ObatRepository {
	return &ObatRepository{}
}

func (r *ObatRepository) FindAll(db *gorm.DB, obat *[]entity.Obat) error {
	return db.Preload("AdminApotek").Find(obat).Error
}
func (r *ObatRepository) FindAllByIdAdminApotek(db *gorm.DB, obat *[]entity.Obat, idAdminApotek int32) error {
	return db.Where("id_admin_apotek = ?", idAdminApotek).Find(obat).Error
}
func (r *ObatRepository) FindById(db *gorm.DB, obat *entity.Obat, id int32) error {
	return db.Where("id = ?", id).First(obat).Error
}
func (r *ObatRepository) FindByIdAndIdAdminApotek(db *gorm.DB, obat *entity.Obat, id int32, idAdminApotek int32) error {
	return db.Where("id = ?", id).Where("id_admin_apotek = ?", idAdminApotek).First(obat).Error
}
func (r *ObatRepository) FindByIdAndLockForUpdate(db *gorm.DB, obat *entity.Obat, id int32) error {
	return db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(obat).Error
}
func (r *ObatRepository) FindByIdAndIdAdminApotekAndLockForUpdate(db *gorm.DB, obat *entity.Obat, id int32, idAdminApotek int32) error {
	return db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).Where("id_admin_apotek = ?", idAdminApotek).First(obat).Error
}
func (r *ObatRepository) FindByIdAdminApotek(db *gorm.DB, obat *entity.Obat, idAdminApotek int32) error {
	return db.Where("id_admin_apotek = ?", idAdminApotek).First(obat).Error
}
