package main

import (
	"context"
	"log"

	"github.com/RandySteven/onboard-be/apps"
	"github.com/RandySteven/onboard-be/configs"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("./files/env/.env"); err != nil {
		log.Println("no .env file loaded, using process environment")
	}
}

func main() {
	configPath, err := configs.ParseFlags()
	if err != nil {
		log.Fatalln(err)
	}

	config, err := configs.NewConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	app, err := apps.NewDBApp(config)
	if err != nil {
		log.Fatalln("Error starting app ", err)
	}

	if err = app.ExecuteMigration(context.Background()); err != nil {
		log.Fatalln("Error executing migration ", err)
	}
	log.Println("success execute migration")
}
