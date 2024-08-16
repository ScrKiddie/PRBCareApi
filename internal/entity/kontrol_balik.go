package entity

type KontrolBalik struct {
	ID             int32  `gorm:"column:id;primaryKey;type:integer;autoIncrement;not null"`
	NoAntrean      int32  `gorm:"column:no_antrean;type:integer;not null"`
	IdPasien       int32  `gorm:"column:id_pasien;type:integer;not null"`
	Pasien         Pasien `gorm:"foreignKey:IdPasien"`
	Keluhan        string `gorm:"column:keluhan;type:text"`
	BeratBadan     int32  `gorm:"column:berat_badan;type:integer"`
	TinggiBadan    int32  `gorm:"column:tinggi_badan;type:integer"`
	TekananDarah   string `gorm:"column:tekanan_darah;type:varchar(20)"`
	DenyutNadi     int32  `gorm:"column:denyut_nadi;type:integer"`
	HasilLab       string `gorm:"column:hasil_lab;type:text"`
	HasilEkg       string `gorm:"column:hasil_ekg;type:text"`
	HasilDiagnosa  string `gorm:"column:hasil_diagnosa;type:text"`
	TanggalKontrol int64  `gorm:"column:tanggal_kontrol;type:bigint;not null"`
	Status         string `gorm:"column:status;type:status_kontrol_balik_enum;not null"`
}

func (KontrolBalik) TableName() string {
	return "kontrol_balik"
}
