package main

import (
	"log"
	"os"

	_ "lamas/migrations"
	"lamas/models"
	"lamas/seeds"
	"lamas/services/caddy"
	"lamas/services/telegram"
	"lamas/services/templates"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	LoadEnv()

	app := pocketbase.New()
	models.RegisterModels(app)

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: false,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		seeds.SeedAll(app)
		return nil
	})

	templates.RegisterService(app)
	caddy.RegisterService(app)
	telegram.RegisterService(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func LoadEnv() {
	log.Println("Loading .env file")
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file")
	}

	if env := os.Getenv("BASE_DOMAIN"); env == "" {
		log.Fatalf("ENV: BASE_DOMAIN must be set")
	}
}
