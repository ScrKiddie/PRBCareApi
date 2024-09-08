package adapter

import (
	"github.com/gofiber/fiber/v3/client"
	"log"
	"prb_care_api/internal/model"
)

type Captcha struct {
	Client *client.Client
}

func NewCaptcha(client *client.Client) *Captcha {
	return &Captcha{Client: client}
}

func (r *Captcha) Verify(request *model.CaptchaRequest) (bool, error) {
	resp, err := r.Client.Post(
		"https://challenges.cloudflare.com/turnstile/v0/siteverify",
		client.Config{
			Body: request,
		},
	)

	if err != nil {
		return false, err
	}

	response := new(model.CaptchaResponse)
	err = resp.JSON(response)
	if err != nil {
		return false, err
	}

	if !response.Success {
		log.Println(string(resp.Body()))
	}
	return response.Success, nil
}
