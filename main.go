package main

import (
	"backend/util"
	"log"
	"net/http"
	"os"

	"backend/controller"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	log.Println("Server started on: server")

	env_err := godotenv.Load(".env")
	util.FailOnError(env_err, ".env Load fail")
	hostIP := os.Getenv("HOST_IP")
	goPort := os.Getenv("GO_SERVER_PORT")

	r := mux.NewRouter()
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://" + hostIP + ":3000", "http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
	})

	handler := corsConfig.Handler(r)

	r.HandleFunc("/yaml", controller.ParseYamlFile).Methods("POST")
	r.HandleFunc("/run", controller.RunWorkflow).Methods("POST")

	http.ListenAndServe(":" + string(goPort), handler)
}
