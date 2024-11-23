package model

type ProlanisResponse struct {
	ID               int32                   `json:"id"`
	IdAdminPuskesmas int32                   `json:"idAdminPuskesmas,omitempty"`
	AdminPuskesmas   *AdminPuskesmasResponse `json:"adminPuskesmas,omitempty"`
	Deskripsi        string                  `json:"deskripsi"`
	WaktuMulai       int64                   `json:"waktuMulai"`
	WaktuSelesai     int64                   `json:"waktuSelesai"`
	Status           string                  `json:"status,omitempty"`
}

type ProlanisSearchRequest struct {
	IdAdminPuskesmas int32  `validate:"omitempty,numeric"`
	Status           string `validate:"omitempty,oneof=aktif selesai"`
}
type ProlanisGetRequest struct {
	ID               int32 `validate:"required,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}
type ProlanisCreateRequest struct {
	Deskripsi        string `json:"deskripsi" validate:"required"`
	WaktuMulai       int64  `json:"waktuMulai" validate:"required,numeric"`
	WaktuSelesai     int64  `json:"waktuSelesai" validate:"required,numeric,gtfield=WaktuMulai"`
	IdAdminPuskesmas int32  `json:"idAdminPuskesmas" validate:"required,numeric"`
}
type ProlanisUpdateRequest struct {
	ID                    int32  `json:"id" validate:"required,numeric"`
	Deskripsi             string `json:"deskripsi" validate:"required"`
	WaktuMulai            int64  `json:"waktuMulai" validate:"required,numeric"`
	WaktuSelesai          int64  `json:"waktuSelesai" validate:"required,numeric,gtfield=WaktuMulai"`
	IdAdminPuskesmas      int32  `json:"idAdminPuskesmas" validate:"required,numeric"`
	CurrentAdminPuskesmas bool   `validate:"omitempty"`
}
type ProlanisDeleteRequest struct {
	ID               int32 `json:"id" validate:"required,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}

type ProlanisSelesaiRequest struct {
	ID               int32 `json:"id" validate:"required,numeric"`
	IdAdminPuskesmas int32 `validate:"omitempty,numeric"`
}
