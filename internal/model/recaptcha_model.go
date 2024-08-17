package model

type RecaptchaResponse struct {
	Success     bool    `json:"success"`
	Score       float64 `json:"score"`
	Action      string  `json:"action"`
	ChallengeTs string  `json:"challenge_ts"`
	Hostname    string  `json:"hostname"`
}
type RecaptchaRequest struct {
	TokenRecaptcha string
	Secret         string
}
