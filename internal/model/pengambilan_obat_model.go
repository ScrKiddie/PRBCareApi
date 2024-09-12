package model

type PengambilanObatResponse struct {
	ID                 int32           `json:"id"`
	Resi               string          `json:"resi,omitempty"`
	IdPasien           int32           `json:"idPasien,omitempty"`
	PasienResponse     *PasienResponse `json:"pasien,omitempty"`
	IdObat             int32           `json:"idObat,omitempty"`
	Obat               *ObatResponse   `json:"obat,omitempty"`
	TanggalPengambilan int64           `json:"tanggalPengambilan"`
	Jumlah             int32           `json:"jumlah"`
	Status             string          `json:"status,omitempty"`
}

type PengambilanObatSearchRequest struct {
	IdPengguna       int32  `validate:"omitempty,numeric"`
	IdAdminPuskesmas int32  `validate:"omitempty,numeric"`
	IdAdminApotek    int32  `validate:"omitempty,numeric"`
	Status           string `validate:"omitempty,oneof=menunggu diambil batal"`
}

type PengambilanObatGetRequest struct {
	ID               int32 `validate:"required,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}
type PengambilanObatCreateRequest struct {
	IdPasien           int32 `json:"idPasien" validate:"required,numeric"`
	IdObat             int32 `json:"idObat" validate:"required,numeric"`
	Jumlah             int32 `json:"jumlah" validate:"required,numeric,gt=0"`
	TanggalPengambilan int64 `json:"tanggalPengambilan" validate:"required,numeric"`
	IdAdminPuskesmas   int32 `validate:"omitempty,numeric"`
}
type PengambilanObatUpdateRequest struct {
	ID                 int32 `json:"id" validate:"required,numeric"`
	IdPasien           int32 `json:"idPasien" validate:"required,numeric"`
	IdObat             int32 `json:"idObat" validate:"required,numeric"`
	Jumlah             int32 `json:"jumlah" validate:"required,numeric,gt=0"`
	TanggalPengambilan int64 `json:"tanggalPengambilan" validate:"required,numeric"`
	IdAdminPuskesmas   int32 `validate:"omitempty,numeric"`
}
type PengambilanObatDeleteRequest struct {
	ID               int32 `json:"id" validate:"required,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}

type PengambilanObatDiambilRequest struct {
	ID            int32 `json:"id" validate:"required,numeric"`
	IdAdminApotek int32 `validate:"omitempty,numeric"`
}

type PengambilanObatBatalRequest struct {
	ID               int32 `json:"id" validate:"required,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}
