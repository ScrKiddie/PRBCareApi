package model

type PasienResponse struct {
	ID               int32                   `json:"id"`
	NoRekamMedis     string                  `json:"noRekamMedis"`
	Pengguna         *PenggunaResponse       `json:"pengguna,omitempty"`
	IdPengguna       int32                   `json:"idPengguna,omitempty"`
	AdminPuskesmas   *AdminPuskesmasResponse `json:"adminPuskesmas,omitempty"`
	IdAdminPuskesmas int32                   `json:"idAdminPuskesmas,omitempty"`
	TanggalDaftar    int64                   `json:"tanggalDaftar"`
	Status           string                  `json:"status,omitempty"`
}

type PasienSearchRequest struct {
	IdPengguna       int32  `validate:"omitempty,numeric"`
	IdAdminPuskesmas int32  `validate:"omitempty,numeric"`
	Status           string `json:"status" validate:"omitempty,oneof=aktif selesai"`
}
type PasienGetRequest struct {
	ID               int32 `json:"id" validate:"required,numeric"`
	IdPengguna       int32 `validate:"omitempty,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}
type PasienCreateRequest struct {
	NoRekamMedis     string `json:"noRekamMedis" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	IdPengguna       int32  `json:"idPengguna" validate:"required,numeric"`
	IdAdminPuskesmas int32  `json:"idAdminPuskesmas" validate:"required,numeric"`
	TanggalDaftar    int64  `json:"tanggalDaftar" validate:"required,numeric"`
}
type PasienUpdateRequest struct {
	ID                    int32  `json:"id" validate:"required,numeric"`
	NoRekamMedis          string `json:"noRekamMedis" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	IdPengguna            int32  `json:"idPengguna" validate:"required,numeric"`
	CurrentAdminPuskesmas bool   `validate:"omitempty"`
	IdAdminPuskesmas      int32  `json:"idAdminPuskesmas" validate:"required,numeric"`
	TanggalDaftar         int64  `json:"tanggalDaftar" validate:"required,numeric"`
}
type PasienDeleteRequest struct {
	ID               int32 `json:"id" validate:"required,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}

type PasienSelesaiRequest struct {
	ID               int32 `json:"id" validate:"required,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}
