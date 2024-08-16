package entity

type Pasien struct {
	ID               int32          `gorm:"column:id;primaryKey;type:integer;autoIncrement;not null"`
	NoRekamMedis     string         `gorm:"column:no_rekam_medis;type:varchar(50);not null"`
	IdPengguna       int32          `gorm:"column:id_pengguna;type:integer;not null"`
	Pengguna         Pengguna       `gorm:"foreignKey:IdPengguna"`
	IdAdminPuskesmas int32          `gorm:"column:id_admin_puskesmas;type:integer;not null"`
	AdminPuskesmas   AdminPuskesmas `gorm:"foreignKey:IdAdminPuskesmas"`
	TanggalDaftar    int64          `gorm:"column:tanggal_daftar;type:bigint;not null"`
	Status           string         `gorm:"column:status;type:status_pasien_enum;not null"`
}

func (Pasien) TableName() string {
	return "pasien"
}
