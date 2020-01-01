package main

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"iwantoask/app"
	_ "iwantoask/docs"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	questionHandler := app.NewQuestionHandler()

	router.HandleFunc("/", questionHandler.ListQuestions).Methods(http.MethodGet)
	router.HandleFunc("/questions", questionHandler.ListQuestions).Methods(http.MethodGet)
	router.HandleFunc("/questions/{id}/{path}", app.ShowQuestion).Methods(http.MethodGet)
	router.HandleFunc("/ask", app.AskQuestion).Methods(http.MethodGet)
	router.HandleFunc("/ask", app.SubmitQuestion).Methods(http.MethodPost)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static/")))

	log.Printf("server started at port: %d", 8080)
	_ = http.ListenAndServe(":8080", router)
}
