package model

type PenggunaResponse struct {
	ID              int32  `json:"id,omitempty"`
	NamaLengkap     string `json:"namaLengkap"`
	Telepon         string `json:"telepon"`
	TeleponKeluarga string `json:"teleponKeluarga"`
	Alamat          string `json:"alamat"`
	Username        string `json:"username,omitempty"`
	Token           string `json:"token,omitempty"`
}

type PenggunaRegisterRequest struct {
	NamaLengkap     string `json:"namaLengkap" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	Telepon         string `json:"telepon" validate:"required,min=10,max=16,not_contain_space"`
	TeleponKeluarga string `json:"teleponKeluarga" validate:"required,min=10,max=16,not_contain_space"`
	Alamat          string `json:"alamat" mod:"normalize_spaces" validate:"required,min=3"`
	Username        string `json:"username" validate:"required,min=6,max=50,not_contain_space"`
	Password        string `json:"password" validate:"required,min=6,max=255,is_password_format,not_contain_space"`
	TokenCaptcha    string `json:"tokenCaptcha" validate:"required,min=100"`
}

type PenggunaLoginRequest struct {
	Username     string `json:"username" validate:"required,min=6,max=50,not_contain_space"`
	Password     string `json:"password" validate:"required,min=6,max=255,is_password_format,not_contain_space"`
	TokenCaptcha string `json:"tokenCaptcha" validate:"required,min=100"`
}
type PenggunaPasswordUpdateRequest struct {
	ID              int32  `validate:"required,numeric"`
	CurrentPassword string `json:"currentPassword" validate:"required,min=6,max=255,is_password_format,not_contain_space"`
	NewPassword     string `json:"newPassword" validate:"required,min=6,max=255,is_password_format,not_contain_space"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=NewPassword"`
}
type PenggunaProfileUpdateRequest struct {
	ID              int32  `json:"id" validate:"required,numeric"`
	NamaLengkap     string `json:"namaLengkap" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	Telepon         string `json:"telepon" validate:"required,min=10,max=16,not_contain_space"`
	TeleponKeluarga string `json:"teleponKeluarga" validate:"required,min=10,max=16,not_contain_space"`
	Alamat          string `json:"alamat" mod:"normalize_spaces" validate:"required,min=3"`
}

type PenggunaVerifyRequest struct {
	ID int32 `validate:"required,numeric"`
}
type PenggunaGetRequest struct {
	ID int32 `json:"id" validate:"required,numeric"`
}
type PenggunaCreateRequest struct {
	NamaLengkap     string `json:"namaLengkap" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	Telepon         string `json:"telepon" validate:"required,min=10,max=16,not_contain_space"`
	TeleponKeluarga string `json:"teleponKeluarga" validate:"required,min=10,max=16,not_contain_space"`
	Alamat          string `json:"alamat" mod:"normalize_spaces" validate:"required,min=3"`
	Username        string `json:"username" validate:"required,min=6,max=50,not_contain_space"`
	Password        string `json:"password" validate:"required,min=6,max=255,is_password_format,not_contain_space"`
}
type PenggunaUpdateRequest struct {
	ID              int32  `json:"id" validate:"required,numeric"`
	NamaLengkap     string `json:"namaLengkap" mod:"normalize_spaces" validate:"required,min=3,max=50"`
	Telepon         string `json:"telepon" validate:"required,min=10,max=16,not_contain_space"`
	TeleponKeluarga string `json:"teleponKeluarga" validate:"required,min=10,max=16,not_contain_space"`
	Alamat          string `json:"alamat" mod:"normalize_spaces" validate:"required,min=3"`
	Username        string `json:"username" validate:"required,min=6,max=50,not_contain_space"`
	Password        string `json:"password" validate:"omitempty,min=6,max=255,is_password_format,not_contain_space"`
}
type PenggunaDeleteRequest struct {
	ID int32 `json:"id" validate:"required,numeric"`
}
type PenggunaTokenPerangkatUpdateRequest struct {
	ID             int32  `validate:"required,numeric"`
	TokenPerangkat string `json:"tokenPerangkat" validate:"required,min=163,max=255"`
}
