package entity

type Artikel struct {
	ID               int32          `gorm:"column:id;primaryKey;type:integer;autoIncrement;not null"`
	IdAdminPuskesmas int32          `gorm:"column:id_admin_puskesmas;type:integer;not null"`
	AdminPuskesmas   AdminPuskesmas `gorm:"foreignKey:IdAdminPuskesmas"`
	Judul            string         `gorm:"column:judul;type:varchar(255);not null"`
	Ringkasan        string         `gorm:"column:ringkasan;type:varchar(1000);not null"`
	Isi              string         `gorm:"column:isi;type:text;not null"`
	TanggalPublikasi int64          `gorm:"column:tanggal_publikasi;type:bigint;not null"`
	Banner           string         `gorm:"column:banner;type:varchar(100);"`
}

func (Artikel) TableName() string {
	return "artikel"
}
