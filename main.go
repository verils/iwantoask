package main

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"iwantoask/app"
	_ "iwantoask/docs"
	"log"
	"net/http"
)

// @Title Iwantoask API
// @Description This is a simple questioning application
// @Version 0.1
//
// @License MIT
//
// @Host localhost:8080
// @BasePath /api
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/questions", app.ListQuestions).Methods(http.MethodGet)
	router.HandleFunc("/api/questions", app.CreateQuestion).Methods(http.MethodPost)
	router.HandleFunc("/api/questions/{id}", app.GetQuestion).Methods(http.MethodGet)

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static/")))

	log.Printf("server started at port: %d", 8080)
	_ = http.ListenAndServe(":8080", router)
}
