package entity

type Pengguna struct {
	ID              int32  `gorm:"column:id;primaryKey;type:integer;autoIncrement;not null"`
	TokenPerangkat  string `gorm:"column:token_perangkat;type:varchar(255)"`
	NamaLengkap     string `gorm:"column:nama_lengkap;type:varchar(100);not null"`
	Telepon         string `gorm:"column:telepon;type:varchar(16);unique;not null"`
	TeleponKeluarga string `gorm:"column:telepon_keluarga;type:varchar(16);not null"`
	Alamat          string `gorm:"column:alamat;type:varchar(500);not null"`
	Username        string `gorm:"column:username;type:varchar(50);unique;not null"`
	Password        string `gorm:"column:password;type:varchar(255);not null"`
}

func (Pengguna) TableName() string {
	return "pengguna"
}
