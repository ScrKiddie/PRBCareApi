package entity

type AdminPuskesmas struct {
	ID               int32  `gorm:"column:id;primaryKey;type:integer;autoIncrement;not null"`
	NamaPuskesmas    string `gorm:"column:nama_puskesmas;type:varchar(100);not null"`
	Telepon          string `gorm:"column:telepon;type:varchar(16);unique;not null"`
	Alamat           string `gorm:"column:alamat;type:varchar(1000);not null"`
	WaktuOperasional string `gorm:"column:waktu_operasional;type:varchar(1000);not null"`
	Username         string `gorm:"column:username;type:varchar(50);unique;not null"`
	Password         string `gorm:"column:password;type:varchar(255);not null"`
}

func (AdminPuskesmas) TableName() string {
	return "admin_puskesmas"
}
