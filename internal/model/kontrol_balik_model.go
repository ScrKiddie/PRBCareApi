package model

type KontrolBalikResponse struct {
	ID             int32           `json:"id"`
	NoAntrean      int32           `json:"noAntrean"`
	IdPasien       int32           `json:"idPasien,omitempty"`
	PasienResponse *PasienResponse `json:"pasien,omitempty"`
	BeratBadan     int32           `json:"beratBadan"`
	TinggiBadan    int32           `json:"tinggiBadan"`
	TekananDarah   string          `json:"tekananDarah"`
	DenyutNadi     int32           `json:"denyutNadi"`
	HasilLab       string          `json:"hasilLab"`
	HasilEkg       string          `json:"hasilEkg"`
	TanggalKontrol int64           `json:"tanggalKontrol"`
	HasilDiagnosa  string          `json:"hasilDiagnosa"`
	Keluhan        string          `json:"keluhan"`
	Status         string          `json:"status,omitempty"`
}

type KontrolBalikSearchRequest struct {
	IdPengguna       int32  `validate:"omitempty,numeric"`
	IdAdminPuskesmas int32  `validate:"omitempty,numeric"`
	Status           string `validate:"omitempty,oneof=menunggu selesai batal"`
}
type KontrolBalikGetRequest struct {
	ID               int32 `validate:"required,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}
type KontrolBalikCreateRequest struct {
	IdPasien         int32 `json:"idPasien" validate:"required,numeric"`
	TanggalKontrol   int64 `json:"tanggalKontrol" validate:"required,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}
type KontrolBalikUpdateRequest struct {
	ID               int32  `json:"id" validate:"required,numeric"`
	NoAntrean        int32  `json:"noAntrean" validate:"required,numeric,gt=0"`
	IdPasien         int32  `json:"idPasien" validate:"required,numeric"`
	TanggalKontrol   int64  `json:"tanggalKontrol" validate:"required,numeric"`
	IdAdminPuskesmas int32  `validate:"omitempty,numeric"`
	BeratBadan       int32  `json:"beratBadan" validate:"numeric,gte=0"`
	TinggiBadan      int32  `json:"tinggiBadan" validate:"numeric,gte=0"`
	TekananDarah     string `json:"tekananDarah" mod:"normalize_spaces" validate:"max=20"`
	DenyutNadi       int32  `json:"denyutNadi" validate:"numeric,gte=0"`
	HasilLab         string `json:"hasilLab"`
	HasilEkg         string `json:"hasilEkg"`
	HasilDiagnosa    string `json:"hasilDiagnosa"`
	Keluhan          string `json:"keluhan"`
}
type KontrolBalikDeleteRequest struct {
	ID               int32 `json:"id" validate:"required,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}

type KontrolBalikSelesaiRequest struct {
	ID               int32 `json:"id" validate:"required,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}

type KontrolBalikBatalRequest struct {
	ID               int32 `json:"id" validate:"required,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}
