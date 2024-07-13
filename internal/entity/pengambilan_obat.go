package entity

type PengambilanObat struct {
	ID                 int         `gorm:"column:id;primaryKey;type:integer;autoIncrement;not null"`
	Resi               string      `gorm:"column:resi;type:varchar(50);not null"`
	IDAdminApotek      int         `gorm:"column:id_admin_apotek;type:integer;not null"`
	AdminApotek        AdminApotek `gorm:"foreignKey:IDAdminApotek"`
	IDPasien           int         `gorm:"column:id_pasien;type:integer;not null"`
	Pasien             Pasien      `gorm:"foreignKey:IDPasien"`
	IDObat             int         `gorm:"column:id_obat;type:integer;not null"`
	Obat               Obat        `gorm:"foreignKey:IDObat"`
	Jumlah             int         `gorm:"column:jumlah;type:integer;not null"`
	TanggalPengambilan int64       `gorm:"column:tanggal_pengambilan;type:bigint;not null"`
	Status             string      `gorm:"column:status;type:status_pengambilan_obat_enum;not null"`
}

func (PengambilanObat) TableName() string {
	return "pengambilan_obat"
}
