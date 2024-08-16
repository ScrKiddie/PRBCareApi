package model

type AdminApotekResponse struct {
	ID               int32  `json:"id,omitempty"`
	NamaApotek       string `json:"namaApotek"`
	Telepon          string `json:"telepon"`
	Alamat           string `json:"alamat"`
	WaktuOperasional string `json:"waktuOperasional"`
	Username         string `json:"username,omitempty"`
	Token            string `json:"token,omitempty"`
}
type AdminApotekLoginRequest struct {
	Username string `json:"username" validate:"required,min=6,max=50,not_contain_space"`
	Password string `json:"password" validate:"required,min=6,max=255,is_password_format,not_contain_space"`
}
type AdminApotekPasswordUpdateRequest struct {
	ID              int32  `validate:"required,numeric"`
	CurrentPassword string `json:"currentPassword" validate:"required,min=6,max=255,is_password_format,not_contain_space"`
	NewPassword     string `json:"newPassword" validate:"required,min=6,max=255,is_password_format,not_contain_space"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=NewPassword"`
}
type AdminApotekProfileUpdateRequest struct {
	ID               int32  `json:"id" validate:"required,numeric"`
	NamaApotek       string `json:"namaApotek" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	Telepon          string `json:"telepon" validate:"required,min=10,max=16,not_contain_space"`
	Alamat           string `json:"alamat" mod:"normalize_spaces" validate:"required,min=3"`
	WaktuOperasional string `json:"waktuOperasional" mod:"normalize_spaces" validate:"required,min=3"`
}
type AdminApotekVerifyRequest struct {
	ID int32 `validate:"required,numeric"`
}
type AdminApotekGetRequest struct {
	ID int32 `json:"id" validate:"required,numeric"`
}
type AdminApotekCreateRequest struct {
	NamaApotek       string `json:"namaApotek" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	Telepon          string `json:"telepon" validate:"required,min=10,max=16,not_contain_space"`
	Alamat           string `json:"alamat" mod:"normalize_spaces" validate:"required,min=3"`
	WaktuOperasional string `json:"waktuOperasional" mod:"normalize_spaces" validate:"required,min=3"`
	Username         string `json:"username" validate:"required,min=6,max=50,not_contain_space"`
	Password         string `json:"password" validate:"required,min=6,max=255,is_password_format,not_contain_space"`
}
type AdminApotekUpdateRequest struct {
	ID               int32  `json:"id" validate:"required,numeric"`
	NamaApotek       string `json:"namaApotek" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	Telepon          string `json:"telepon" validate:"required,min=10,max=16,not_contain_space"`
	Alamat           string `json:"alamat" mod:"normalize_spaces" validate:"required,min=3"`
	WaktuOperasional string `json:"waktuOperasional" mod:"normalize_spaces" validate:"required,min=3"`
	Username         string `json:"username" validate:"required,min=6,max=50,not_contain_space"`
	Password         string `json:"password" validate:"omitempty,min=6,max=255,is_password_format,not_contain_space"`
}
type AdminApotekDeleteRequest struct {
	ID int32 `json:"id" validate:"required,numeric"`
}
