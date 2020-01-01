package app

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"time"
)

type UserModel struct {
	Username sql.NullString
	Name     sql.NullString
}

func (userModel *UserModel) GetUsernameOrDefault(value string) string {
	if userModel.Username.Valid {
		return userModel.Username.String
	}
	return value
}

func (userModel *UserModel) GetNameOrDefault(value string) string {
	if userModel.Name.Valid {
		return userModel.Name.String
	}
	return value
}

type User struct {
	Username string
	Name     string
}

type Question struct {
	Id      int
	Title   string
	Detail  string
	Path    string
	Url     string
	AskedAt time.Time
	AskedBy User
}

type QuestionsView struct {
	Questions []Question
}

type QuestionHandler struct {
}

func NewQuestionHandler() *QuestionHandler {
	return &QuestionHandler{}
}

func (handler *QuestionHandler) ListQuestions(writer http.ResponseWriter, request *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(docker.local)/iwantoask?parseTime=true")
	if err != nil {
		http.Error(writer, fmt.Sprintf("cannot connect to database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT q.id, q.path, q.title, q.detail, q.asked_at, q.asked_by, u.name FROM questions q LEFT JOIN users u ON q.asked_by = u.username")
	if err != nil {
		http.Error(writer, fmt.Sprintf("cannot fetch questions: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	questions := make([]Question, 0)
	for rows.Next() {
		var question Question
		var userModel UserModel

		err = rows.Scan(&question.Id, &question.Path, &question.Title, &question.Detail, &question.AskedAt, &userModel.Username, &userModel.Name)
		if err != nil {
			http.Error(writer, fmt.Sprintf("scan error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		question.AskedBy = User{
			Username: userModel.GetUsernameOrDefault("anonymous"),
			Name:     userModel.GetNameOrDefault("anonymous"),
		}
		questions = append(questions, question)
	}

	for i, question := range questions {
		questions[i].Url = fmt.Sprintf("http://%s/q/%d/%s", request.Host, question.Id, question.Path)
	}

	tmpl := template.Must(template.ParseFiles("template/questions.html"))
	_ = tmpl.Execute(writer, questions)
}

func ShowQuestion(writer http.ResponseWriter, request *http.Request) {

}

func AskQuestion(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/ask.html"))
	_ = tmpl.Execute(writer, nil)
}

func SubmitQuestion(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Post Form: %s", request.PostFormValue("title"))
}
