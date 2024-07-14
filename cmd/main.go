package main

import (
	"log"
	"prbcare_be/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	app := config.NewFiber()
	db := config.NewDatabase(viperConfig)
	validator := config.NewValidator()
	mold := config.NewMold()
	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Validate: validator,
		Config:   viperConfig,
		Modifier: mold,
	})
	err := app.Listen("localhost:" + viperConfig.GetString("web.port"))
	if err != nil {
		log.Fatalln(err)
	}
}
