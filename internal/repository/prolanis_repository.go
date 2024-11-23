package repository

import (
	"gorm.io/gorm"
	"prb_care_api/internal/entity"
)

type ProlanisRepository struct {
	Repository[entity.Prolanis]
}

func NewProlanisRepository() *ProlanisRepository {
	return &ProlanisRepository{}
}

func (r *ProlanisRepository) Search(db *gorm.DB, prolanis *[]entity.Prolanis, status string, idAdminPuskesmas int32) error {
	query := db
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if idAdminPuskesmas != 0 {
		query = query.Where("id_admin_puskesmas = ?", idAdminPuskesmas)
	}
	return query.Preload("AdminPuskesmas").Find(&prolanis).Error
}

func (r *ProlanisRepository) FindAll(db *gorm.DB, prolanis *[]entity.Prolanis) error {
	return db.Preload("prolanis").Find(prolanis).Error
}

func (r *ProlanisRepository) FindAllByIdAdminPuskesmas(db *gorm.DB, prolanis *[]entity.Prolanis, idAdminPuskesmas int32) error {
	return db.Where("id_admin_puskesmas = ?", idAdminPuskesmas).Find(prolanis).Error
}

func (r *ProlanisRepository) FindByIdAndStatus(db *gorm.DB, prolanis *entity.Prolanis, id int32, status string) error {
	return db.Where("id = ?", id).Where("status = ?", status).First(prolanis).Error
}
func (r *ProlanisRepository) FindByIdAndIdAdminPuskesmasAndStatus(db *gorm.DB, prolanis *entity.Prolanis, id int32, idAdminPuskesmas int32, status string) error {
	return db.Where("id = ?", id).Where("id_admin_puskesmas = ?", idAdminPuskesmas).Where("status = ?", status).First(prolanis).Error
}
func (r *ProlanisRepository) FindByIdAdminPuskesmas(db *gorm.DB, prolanis *entity.Prolanis, idAdminPuskesmas int32) error {
	return db.Where("id_admin_puskesmas = ?", idAdminPuskesmas).First(prolanis).Error
}
