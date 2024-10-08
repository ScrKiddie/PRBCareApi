package repository

import (
	"gorm.io/gorm"
	"prb_care_api/internal/entity"
)

type AdminPuskesmasRepository struct {
	Repository[entity.AdminPuskesmas]
}

func NewAdminPuskesmasRepository() *AdminPuskesmasRepository {
	return &AdminPuskesmasRepository{}
}

func (r *AdminPuskesmasRepository) FindByUsername(db *gorm.DB, adminPuskesmas *entity.AdminPuskesmas, username string) error {
	return db.Where("username = ?", username).First(adminPuskesmas).Error
}
func (r *AdminPuskesmasRepository) CountByUsername(db *gorm.DB, username any) (int64, error) {
	var count int64
	if err := db.Model(&entity.AdminPuskesmas{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
func (r *AdminPuskesmasRepository) CountByTelepon(db *gorm.DB, telepon any) (int64, error) {
	var count int64
	if err := db.Model(&entity.AdminPuskesmas{}).Where("telepon = ?", telepon).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
func (r *AdminPuskesmasRepository) FindAll(db *gorm.DB, adminPuskesmas *[]entity.AdminPuskesmas) error {
	return db.Find(adminPuskesmas).Error
}
func (r *AdminPuskesmasRepository) FindById(db *gorm.DB, adminPuskesmas *entity.AdminPuskesmas, id int32) error {
	return db.Where("id = ?", id).First(adminPuskesmas).Error
}
