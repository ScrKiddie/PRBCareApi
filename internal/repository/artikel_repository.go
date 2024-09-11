package repository

import (
	"gorm.io/gorm"
	"prb_care_api/internal/entity"
)

type ArtikelRepository struct {
	Repository[entity.Artikel]
}

func NewArtikelRepository() *ArtikelRepository {
	return &ArtikelRepository{}
}

func (r *ArtikelRepository) Search(db *gorm.DB, artikel *[]entity.Artikel, idAdminPuskesmas int32) error {
	query := db
	if idAdminPuskesmas != 0 {
		query = query.Where("id_admin_puskesmas = ?", idAdminPuskesmas)
	}
	return query.Preload("AdminPuskesmas").Find(artikel).Error
}

func (r *ArtikelRepository) FindById(db *gorm.DB, artikel *entity.Artikel, id int32) error {
	return db.Where("id = ?", id).Preload("AdminPuskesmas").First(artikel).Error
}
func (r *ArtikelRepository) FindByIdAndIdAdminPuskesmas(db *gorm.DB, artikel *entity.Artikel, idAdminPuskesmas int32, id int32) error {
	return db.Where("id = ?", id).Where("id_admin_puskesmas = ?", idAdminPuskesmas).Preload("AdminPuskesmas").First(artikel).Error
}
