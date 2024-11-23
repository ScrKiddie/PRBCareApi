package entity

type Prolanis struct {
	ID               int32          `gorm:"column:id;primaryKey;type:integer;autoIncrement;not null"`
	IdAdminPuskesmas int32          `gorm:"column:id_admin_puskesmas;type:integer;not null"`
	AdminPuskesmas   AdminPuskesmas `gorm:"foreignKey:IdAdminPuskesmas"`
	Deskripsi        string         `gorm:"column:deskripsi;type:text"`
	WaktuMulai       int64          `gorm:"column:waktu_mulai;type:bigint;not null"`
	WaktuSelesai     int64          `gorm:"column:waktu_selesai;type:bigint;not null"`
	Status           string         `gorm:"column:status;type:status_prolanis_enum;not null"`
}

func (Prolanis) TableName() string {
	return "prolanis"
}
