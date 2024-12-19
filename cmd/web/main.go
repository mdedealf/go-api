package main

import (
	"fmt"
	"github.com/mdedealf/go-api/internal/config"
)

func main() {
	//app := fiber.New()
	//app.Get("/ping", func(c *fiber.Ctx) error {
	//	return c.SendString("pong")
	//})
	//
	//app.Listen(":8081")

	viperConfig := config.LoadConfig()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	app := config.NewFiber(viperConfig)

	cfg := &config.AppConfig{
		DB:     db,
		App:    app,
		Log:    log,
		Config: viperConfig,
	}

	cfg.Run()

	webPort := viperConfig.GetInt("APP_PORT")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
