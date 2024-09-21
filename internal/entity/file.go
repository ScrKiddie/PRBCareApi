package entity

type File struct {
	ID        int32   `gorm:"column:id;primaryKey;type:integer;autoIncrement;not null"`
	IdArtikel int32   `gorm:"column:id_artikel;type:integer;not null"`
	Artikel   Artikel `gorm:"foreignKey:IdArtikel"`
	File      string  `gorm:"column:file;type:varchar(100);"`
}

func (File) TableName() string {
	return "file"
}
