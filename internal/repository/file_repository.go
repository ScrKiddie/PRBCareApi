package repository

import (
	"gorm.io/gorm"
	"prb_care_api/internal/entity"
)

type FileRepository struct {
	Repository[entity.File]
}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (r *FileRepository) SearchByIdArtikel(db *gorm.DB, file *[]entity.File, idArtikel int32) error {
	return db.Where("id_artikel = ?", idArtikel).Find(file).Error
}
