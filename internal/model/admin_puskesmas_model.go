package model

type AdminPuskesmasResponse struct {
	ID               int32  `json:"id,omitempty"`
	NamaPuskesmas    string `json:"namaPuskesmas"`
	Telepon          string `json:"telepon"`
	Alamat           string `json:"alamat"`
	Username         string `json:"username,omitempty"`
	WaktuOperasional string `json:"waktuOperasional"`
	Token            string `json:"token,omitempty"`
}
type AdminPuskesmasLoginRequest struct {
	Username     string `json:"username" validate:"required,min=6,max=50,not_contain_space"`
	Password     string `json:"password" validate:"required,min=6,max=255,is_password_format,not_contain_space"`
	TokenCaptcha string `json:"tokenCaptcha" validate:"required,min=100"`
}
type AdminPuskesmasPasswordUpdateRequest struct {
	ID              int32  `validate:"required,numeric"`
	CurrentPassword string `json:"currentPassword" validate:"required,min=6,max=255,is_password_format,not_contain_space"`
	NewPassword     string `json:"newPassword" validate:"required,min=6,max=255,is_password_format,not_contain_space"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=NewPassword"`
}
type AdminPuskesmasProfileUpdateRequest struct {
	ID               int32  `json:"id" validate:"required,numeric"`
	NamaPuskesmas    string `json:"namaPuskesmas" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	Telepon          string `json:"telepon" validate:"required,min=10,max=16,not_contain_space"`
	Alamat           string `json:"alamat" mod:"normalize_spaces" validate:"required,min=3,max=1000"`
	WaktuOperasional string `json:"waktuOperasional" mod:"normalize_spaces" validate:"required,min=3,max=1000"`
}
type AdminPuskesmasVerifyRequest struct {
	ID int32 `validate:"required,numeric"`
}
type AdminPuskesmasGetRequest struct {
	ID int32 `validate:"required,numeric"`
}
type AdminPuskesmasCreateRequest struct {
	NamaPuskesmas    string `json:"namaPuskesmas" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	Telepon          string `json:"telepon" validate:"required,min=10,max=16,not_contain_space"`
	Alamat           string `json:"alamat" mod:"normalize_spaces" validate:"required,min=3,max=1000"`
	WaktuOperasional string `json:"waktuOperasional" mod:"normalize_spaces" validate:"required,min=3,max=1000"`
	Username         string `json:"username" validate:"required,min=6,max=50,not_contain_space"`
	Password         string `json:"password" validate:"required,min=6,max=255,is_password_format,not_contain_space"`
}
type AdminPuskesmasUpdateRequest struct {
	ID               int32  `json:"id" validate:"required,numeric"`
	NamaPuskesmas    string `json:"namaPuskesmas" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	Telepon          string `json:"telepon" validate:"required,min=10,max=16,not_contain_space"`
	Alamat           string `json:"alamat" mod:"normalize_spaces" validate:"required,min=3,max=1000"`
	WaktuOperasional string `json:"waktuOperasional" mod:"normalize_spaces" validate:"required,min=3,max=1000"`
	Username         string `json:"username" validate:"required,min=6,max=50,not_contain_space"`
	Password         string `json:"password" validate:"omitempty,min=6,max=255,is_password_format,not_contain_space"`
}
type AdminPuskesmasDeleteRequest struct {
	ID int32 `json:"id" validate:"required,numeric"`
}
