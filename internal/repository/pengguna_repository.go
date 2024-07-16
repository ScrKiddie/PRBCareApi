package repository

import (
	"gorm.io/gorm"
	"prbcare_be/internal/entity"
)

type PenggunaRepository struct {
}

func NewPenggunaRepository() *PenggunaRepository {
	return &PenggunaRepository{}
}
func (r *PenggunaRepository) FindByUsername(db *gorm.DB, pengguna *entity.Pengguna, username string) error {
	return db.Where("username = ?", username).First(pengguna).Error
}
func (r *PenggunaRepository) FindById(db *gorm.DB, pengguna *entity.Pengguna, id int) error {
	return db.Where("id = ?", id).First(pengguna).Error
}
func (r *PenggunaRepository) Update(db *gorm.DB, pengguna *entity.Pengguna) error {
	return db.Save(pengguna).Error
}
func (r *PenggunaRepository) Delete(db *gorm.DB, pengguna *entity.Pengguna) error {
	return db.Delete(pengguna).Error
}
func (r *PenggunaRepository) Create(db *gorm.DB, pengguna *entity.Pengguna) error {
	return db.Create(pengguna).Error
}
func (r *PenggunaRepository) CountByUsername(db *gorm.DB, username any) (int64, error) {
	var count int64
	if err := db.Model(&entity.Pengguna{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
func (r *PenggunaRepository) CountByTelepon(db *gorm.DB, telepon any) (int64, error) {
	var count int64
	if err := db.Model(&entity.Pengguna{}).Where("telepon = ?", telepon).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
func (r *PenggunaRepository) FindAll(db *gorm.DB, pengguna *[]entity.Pengguna) error {
	return db.Find(pengguna).Error
}
