package entity

type AdminApotek struct {
	ID         int    `gorm:"column:id;primaryKey;type:integer;autoIncrement;not null"`
	NamaApotek string `gorm:"column:nama_apotek;type:varchar(100);not null"`
	Phone      string `gorm:"column:phone;type:varchar(15);unique;not null"`
	Address    string `gorm:"column:address;type:text;not null"`
	Username   string `gorm:"column:username;type:varchar(50);unique;not null"`
	Password   string `gorm:"column:password;type:varchar(255);not null"`
}

func (AdminApotek) TableName() string {
	return "admin_apotek"
}
