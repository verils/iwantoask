package main

import (
	"github.com/gorilla/mux"
	"iwantoask/app"
	"log"
	"net/http"
	"os"
)

func main() {
	if os.Getenv(app.MysqlHost) == "" {
		_ = os.Setenv(app.MysqlHost, "docker.local")
	}
	if os.Getenv(app.MysqlUsername) == "" {
		_ = os.Setenv(app.MysqlUsername, "root")
	}
	if os.Getenv(app.MysqlPassword) == "" {
		_ = os.Setenv(app.MysqlPassword, "")
	}

	router := mux.NewRouter()

	router.HandleFunc("/", app.ListQuestions).Methods(http.MethodGet)
	router.HandleFunc("/questions", app.ListQuestions).Methods(http.MethodGet)
	router.HandleFunc("/ask", app.AskQuestion).Methods(http.MethodGet)
	router.HandleFunc("/ask", app.SubmitQuestion).Methods(http.MethodPost)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static/")))

	log.Printf("server started at port: %d", 8080)
	_ = http.ListenAndServe(":8080", router)
}
