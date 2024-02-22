package api

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var hostIP string
var goPort string
var pythonPort string

func init() {
	hostIP = os.Getenv("HOST_IP")
	goPort = os.Getenv("GO_SERVER_PORT")
	pythonPort = os.Getenv("PYTHON_SERVER_PORT")
}

func InitRouter() {
	r := mux.NewRouter()
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://" + hostIP + ":3000", "http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
	})

	handler := corsConfig.Handler(r)

	r.HandleFunc("/api/v1/yaml", ParseYamlFile).Methods("POST")

	r.HandleFunc("/api/v1/info", GetWorkflowInfos).Methods("GET")
	r.HandleFunc("/api/v1/run", RunWorkflow).Methods("POST")

	http.ListenAndServe(":"+string(goPort), handler)
}
