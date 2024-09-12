package model

type ObatResponse struct {
	ID            int32                `json:"id"`
	IdAdminApotek int32                `json:"idAdminApotek,omitempty"`
	AdminApotek   *AdminApotekResponse `json:"adminApotek,omitempty"`
	NamaObat      string               `json:"namaObat"`
	Jumlah        int32                `json:"jumlah"`
}

type ObatListRequest struct {
	IdAdminApotek int32 `validate:"omitempty,numeric"`
}
type ObatGetRequest struct {
	ID            int32 `validate:"required,numeric"`
	IdAdminApotek int32 `validate:"omitempty,numeric"`
}
type ObatCreateRequest struct {
	NamaObat      string `json:"namaObat" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	Jumlah        int32  `json:"jumlah" validate:"required,numeric,gt=0"`
	IdAdminApotek int32  `validate:"required,numeric"`
}
type ObatUpdateRequest struct {
	ID                 int32  `json:"id" validate:"required,numeric"`
	NamaObat           string `json:"namaObat" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	Jumlah             int32  `json:"jumlah" validate:"numeric,gte=0"`
	CurrentAdminApotek bool   `validate:"omitempty"`
	IdAdminApotek      int32  `json:"idAdminApotek" validate:"required,numeric"`
}
type ObatDeleteRequest struct {
	ID            int32 `json:"id" validate:"required,numeric"`
	IdAdminApotek int32 `validate:"omitempty,numeric"`
}
