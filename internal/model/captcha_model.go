package model

type CaptchaResponse struct {
	Success bool `json:"success"`
}
type CaptchaRequest struct {
	Secret       string `json:"secret"`
	TokenCaptcha string `json:"response"`
}
