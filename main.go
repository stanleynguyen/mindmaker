package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stanleynguyen/mindmaker/mindmaker"
	"github.com/stanleynguyen/mindmaker/persistence/postgres"
)

func main() {
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
	}

	// database, err := redis.NewInstance(os.Getenv("DB"))
	database, err := postgres.NewInstance(os.Getenv("DB"))
	if err != nil {
		log.Fatal(err)
	}
	err = mindmaker.Initialize(mindmaker.Config{
		Token:         os.Getenv("BOT_TOKEN"),
		SSL:           false,
		WebhookAddr:   os.Getenv("WEBHOOK_URL"),
		ListeningPath: "/",
	}, database)
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
