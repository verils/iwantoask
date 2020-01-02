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
		defaultMysqlHost := "docker.local"
		_ = os.Setenv(app.MysqlHost, defaultMysqlHost)
		log.Printf("[WARN] missing env %s, set default to '%s'", app.MysqlHost, defaultMysqlHost)
	}
	if os.Getenv(app.MysqlUsername) == "" {
		defaultMysqlUsername := "root"
		_ = os.Setenv(app.MysqlUsername, defaultMysqlUsername)
		log.Printf("[WARN] missing env %s, set default to '%s'", app.MysqlUsername, defaultMysqlUsername)
	}
	if os.Getenv(app.MysqlPassword) == "" {
		defaultMysqlPassword := ""
		_ = os.Setenv(app.MysqlPassword, defaultMysqlPassword)
		log.Printf("[WARN] missing env %s, set default to '%s'", app.MysqlPassword, defaultMysqlPassword)
	}

	router := mux.NewRouter()

	router.HandleFunc("/iwantoask/", app.ListQuestions).Methods(http.MethodGet)
	router.HandleFunc("/iwantoask/questions", app.ListQuestions).Methods(http.MethodGet)
	router.HandleFunc("/iwantoask/ask", app.AskQuestion).Methods(http.MethodGet)
	router.HandleFunc("/iwantoask/ask", app.SubmitQuestion).Methods(http.MethodPost)

	router.PathPrefix("/iwantoask/").Handler(http.StripPrefix("/iwantoask/", http.FileServer(http.Dir("static/"))))

	log.Printf("[INFO] server started at port: %d", 8080)
	_ = http.ListenAndServe(":8080", router)
}
