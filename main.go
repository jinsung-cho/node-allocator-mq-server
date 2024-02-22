package main

import (
	"log"
	"server/internal/api"
	"server/pkg/handler"

	"github.com/joho/godotenv"
)

func main() {
	env_err := godotenv.Load(".env")
	handler.CheckErrorAndPanic(env_err, ".env Load fail")

	api.InitRouter()

	log.Println("Server started on: server")
}
