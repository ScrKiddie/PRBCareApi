package adapter

import (
	"fmt"
	"github.com/gofiber/fiber/v3/client"
	"log"
	"prb_care_api/internal/model"
)

type Recaptcha struct {
	Client *client.Client
}

func NewRecaptcha(client *client.Client) *Recaptcha {
	return &Recaptcha{Client: client}
}

func (r *Recaptcha) Verify(request *model.RecaptchaRequest) (bool, error) {
	resp, err := r.Client.Post(
		fmt.Sprintf("https://www.google.com/recaptcha/api/siteverify?secret=%s&response=%s",
			request.Secret, request.TokenRecaptcha), client.Config{})

	if err != nil {
		return false, err
	}

	response := new(model.RecaptchaResponse)
	err = resp.JSON(response)
	if err != nil {
		return false, err
	}

	if !response.Success {
		log.Println(string(resp.Body()))
	}
	return response.Success, nil
}
