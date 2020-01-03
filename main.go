package main

import (
	"github.com/gorilla/mux"
	"iwantoask/app"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Printf("[INFO] iwantoask (iwantoask-%s) starting...", app.Version)
	log.Printf("[INFO] iwantoask web base path: '%s'", app.BasePath)

	initializeEnvironment()

	router := mux.NewRouter()

	router.HandleFunc(app.BasePathPrefix("/"), app.ListQuestions).Methods(http.MethodGet)
	router.HandleFunc(app.BasePathPrefix("/questions"), app.ListQuestions).Methods(http.MethodGet)
	router.HandleFunc(app.BasePathPrefix("/ask"), app.AskQuestion).Methods(http.MethodGet)
	router.HandleFunc(app.BasePathPrefix("/ask"), app.SubmitQuestion).Methods(http.MethodPost)

	router.PathPrefix(app.BasePathPrefix("/")).Handler(http.StripPrefix(app.BasePathPrefix("/"), http.FileServer(http.Dir("static/"))))

	log.Printf("[INFO] server started at port: %d", 8080)
	_ = http.ListenAndServe(":8080", router)
}

func initializeEnvironment() {
	if os.Getenv(app.EnvMysqlHost) == "" {
		defaultMysqlHost := "docker.local"
		_ = os.Setenv(app.EnvMysqlHost, defaultMysqlHost)
		log.Printf("[WARN] missing env %s, set default to '%s'", app.EnvMysqlHost, defaultMysqlHost)
	}
	if os.Getenv(app.EnvMysqlUsername) == "" {
		defaultMysqlUsername := "root"
		_ = os.Setenv(app.EnvMysqlUsername, defaultMysqlUsername)
		log.Printf("[WARN] missing env %s, set default to '%s'", app.EnvMysqlUsername, defaultMysqlUsername)
	}
	if os.Getenv(app.EnvMysqlPassword) == "" {
		defaultMysqlPassword := ""
		_ = os.Setenv(app.EnvMysqlPassword, defaultMysqlPassword)
		log.Printf("[WARN] missing env %s, set default to '%s'", app.EnvMysqlPassword, defaultMysqlPassword)
	}
}
