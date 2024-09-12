package model

type ArtikelResponse struct {
	ID               int32                   `json:"id"`
	AdminPuskesmas   *AdminPuskesmasResponse `json:"adminPuskesmas,omitempty"`
	IdAdminPuskesmas int32                   `json:"idAdminPuskesmas,omitempty"`
	Judul            string                  `json:"judul"`
	Ringkasan        string                  `json:"ringkasan,omitempty"`
	Isi              string                  `json:"isi,omitempty"`
	TanggalPublikasi int64                   `json:"tanggalPublikasi"`
}

type ArtikelGetRequest struct {
	ID int32 `validate:"required,numeric"`
}
type ArtikelSearchRequest struct {
	IdAdminPuskesmas int32 `validate:"omitempty,numeric,gte=0"`
}
type ArtikelCreateRequest struct {
	Judul            string `json:"judul" mod:"normalize_spaces" validate:"required,max=255"`
	Ringkasan        string `json:"ringkasan" mod:"normalize_spaces" validate:"required,max=1000"`
	Isi              string `json:"isi" validate:"required"`
	IdAdminPuskesmas int32  `json:"idAdminPuskesmas" validate:"required,numeric"`
}

type ArtikelUpdateRequest struct {
	ID                    int32  `json:"id" validate:"required,numeric"`
	Judul                 string `json:"judul" mod:"normalize_spaces" validate:"required,max=255"`
	Ringkasan             string `json:"ringkasan" mod:"normalize_spaces" validate:"required,max=1000"`
	Isi                   string `json:"isi" validate:"required"`
	IdAdminPuskesmas      int32  `json:"idAdminPuskesmas" validate:"required,numeric"`
	CurrentAdminPuskesmas bool
}

type ArtikelDeleteRequest struct {
	ID               int32 `json:"id" validate:"required,numeric"`
	IdAdminPuskesmas int32 `json:"idAdminPuskesmas" validate:"required,numeric"`
}
