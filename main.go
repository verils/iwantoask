package main

import (
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/verils/iwantoask/app"
	"log"
	"net/http"
	"time"
)

const Version = "0.1.0"

func main() {
	log.Printf("[INFO] iwantoask (iwantoask-%s) starting...", Version)
	log.Printf("[INFO] iwantoask base path: '%s'", app.BasePath)

	db := initDB()
	defer db.Close()

	handler := app.NewQuestionHandler(db)

	router := mux.NewRouter()

	router.HandleFunc(app.PrefixBasePath("/"), app.CookieAuth(handler.ListQuestions)).Methods(http.MethodGet)
	router.HandleFunc(app.PrefixBasePath("/questions"), app.CookieAuth(handler.ListQuestions)).Methods(http.MethodGet)
	router.HandleFunc(app.PrefixBasePath("/questions.json"), app.CookieAuth(handler.ListQuestionsJson)).Methods(http.MethodGet)
	router.HandleFunc(app.PrefixBasePath("/ask"), app.CookieAuth(handler.AskQuestion)).Methods(http.MethodGet)
	router.HandleFunc(app.PrefixBasePath("/ask"), app.CookieAuth(handler.SubmitQuestion)).Methods(http.MethodPost)

	router.PathPrefix(app.PrefixBasePath("/")).Handler(http.StripPrefix(app.PrefixBasePath("/"), http.FileServer(http.Dir("static/"))))

	log.Printf("[INFO] server started at port: %d", 8080)
	_ = http.ListenAndServe(":8080", router)
}

func initDB() *bolt.DB {
	db, err := bolt.Open("iwantoask.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("[ERROR] failed to open database: %s", err.Error())
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists([]byte(app.BucketQuestions))
		return e
	})
	if err != nil {
		log.Fatalf("[ERROR] failed to create bucket: %s", err.Error())
	}

	return db
}
