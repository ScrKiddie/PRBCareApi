package main

import (
	"log"
	"prb_care_api/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	app := config.NewFiber()
	db := config.NewDatabase(viperConfig)
	validator := config.NewValidator()
	mold := config.NewMold()
	client := config.NewClient()
	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Validate: validator,
		Config:   viperConfig,
		Modifier: mold,
		Client:   client,
	})
	err := app.Listen("0.0.0.0:" + viperConfig.GetString("web.port"))
	if err != nil {
		log.Fatalln(err)
	}
}
