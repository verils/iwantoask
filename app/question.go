package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	haikunator "github.com/atrox/haikunatorgo/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Question struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Path      string    `json:"path"`
	Url       string    `json:"url"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// ListQuestions
//
// @Summary ListQuestions
// @Description list all questions
// @Tags Questions API
//
// @Router /questions [get]
// @Accept json
// @Produce json
//
// @Success 200 {array} Question "the result question data"
// @Success 500 {object} ErrorResult "error messages"
func ListQuestions(writer http.ResponseWriter, request *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(docker.local)/iwantoask?parseTime=true")
	if err != nil {
		bytes, _ := json.Marshal(ErrorResult{Message: "cannot connect to database", Details: err.Error()})
		http.Error(writer, string(bytes), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, path, title, content, created_by, created_at FROM questions")
	if err != nil {
		bytes, _ := json.Marshal(ErrorResult{Message: "cannot fetch questions", Details: err.Error()})
		http.Error(writer, string(bytes), http.StatusInternalServerError)
		return
	}

	questions := make([]Question, 0)
	for rows.Next() {
		var question Question
		err = rows.Scan(&question.Id, &question.Path, &question.Title, &question.Content, &question.CreatedBy, &question.CreatedAt)
		if err != nil {
			bytes, _ := json.Marshal(ErrorResult{Message: "scan error", Details: err.Error()})
			http.Error(writer, string(bytes), http.StatusInternalServerError)
			return
		}
		questions = append(questions, question)
	}

	schema := "http"
	if request.TLS != nil {
		schema = "https"
	}
	for i, question := range questions {
		questions[i].Url = fmt.Sprintf("%s://%s/q/%d/%s", schema, request.Host, question.Id, question.Path)
	}

	_ = json.NewEncoder(writer).Encode(&questions)
}

// CreateQuestion
//
// @Summary CreateQuestion
// @Description create a question
// @Tags Questions API
//
// @Router /questions [post]
// @Accept json
// @Produce json
// @Param question body Question true "question content"
// @Success 201 {object} Question "the created question data"
// @Success 500 {object} ErrorResult "error messages"
func CreateQuestion(writer http.ResponseWriter, request *http.Request) {
	var question Question
	_ = json.NewDecoder(request.Body).Decode(&question)

	title := question.Title
	question.Path = strings.ReplaceAll(strings.ToLower(title), " ", "+")

	haikunate := haikunator.New()
	haikunate.Delimiter = "_"
	haikunate.TokenLength = 0
	username := haikunate.Haikunate()

	question.CreatedBy = username
	question.CreatedAt = time.Now()

	db, err := sql.Open("mysql", "root:@tcp(docker.local)/iwantoask?parseTime=true")
	if err != nil {
		bytes, _ := json.Marshal(ErrorResult{Message: "cannot connect to database", Details: err.Error()})
		http.Error(writer, string(bytes), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO questions (path, title, content, created_by, created_at) VALUES (?, ?, ?, ?, ?)",
		question.Path, question.Title, question.Content, question.CreatedBy, question.CreatedAt)
	if err != nil {
		bytes, _ := json.Marshal(ErrorResult{Message: "an error happened when accessing the resources", Details: err.Error()})
		http.Error(writer, string(bytes), http.StatusInternalServerError)
		return
	}
	lastInsertId, _ := result.LastInsertId()
	question.Id = int(lastInsertId)

	schema := "http"
	if request.TLS != nil {
		schema = "https"
	}
	question.Url = fmt.Sprintf("%s://%s/q/%d/%s", schema, request.Host, question.Id, question.Path)

	writer.Header().Set("Location", question.Url)
	writer.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(writer).Encode(&question)
}

// GetQuestion
//
// @Summary GetQuestion
// @Description get one question
// @Tags Questions API
//
// @Router /questions/{id} [get]
// @Accept json
// @Produce json
// @Param id path int true "question id"
// @Success 200 {array} Question "the result question data"
// @Success 500 {object} ErrorResult "error messages"
func GetQuestion(writer http.ResponseWriter, request *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(docker.local)/iwantoask?parseTime=true")
	if err != nil {
		bytes, _ := json.Marshal(ErrorResult{Message: "cannot connect to database", Details: err.Error()})
		http.Error(writer, string(bytes), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	vars := mux.Vars(request)

	paramId := vars["id"]
	id, err := strconv.ParseInt(paramId, 10, 32)
	if err != nil {
		bytes, _ := json.Marshal(ErrorResult{Message: "question resource not found, id: " + paramId, Details: err.Error()})
		http.Error(writer, string(bytes), http.StatusNotFound)
		return
	}

	row := db.QueryRow("SELECT id, path, title, content, created_by, created_at FROM questions WHERE id = ?", id)

	var question Question
	err = row.Scan(&question.Id, &question.Path, &question.Title, &question.Content, &question.CreatedBy, &question.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			bytes, _ := json.Marshal(ErrorResult{Message: "question resource not found, id: " + paramId, Details: err.Error()})
			http.Error(writer, string(bytes), http.StatusNotFound)
			return
		} else {
			bytes, _ := json.Marshal(ErrorResult{Message: "cannot fetch question: " + paramId, Details: err.Error()})
			http.Error(writer, string(bytes), http.StatusInternalServerError)
			return
		}
	}

	schema := "http"
	if request.TLS != nil {
		schema = "https"
	}
	question.Url = fmt.Sprintf("%s://%s/q/%d/%s", schema, request.Host, question.Id, question.Path)

	_ = json.NewEncoder(writer).Encode(&question)
}
