package config

import (
	"github.com/gofiber/fiber/v3/client"
	"time"
)

func NewClient() *client.Client {
	c := client.New()
	c.SetTimeout(10 * time.Second)
	return c
}
