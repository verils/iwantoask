package app

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const MysqlHost = "MYSQL_HOST"
const MysqlUsername = "MYSQL_USERNAME"
const MysqlPassword = "MYSQL_PASSWORD"

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
	Since   string
}

type QuestionsView struct {
	Questions         []Question
	SortByRecently    bool
	SortByInteresting bool
}

func ListQuestions(writer http.ResponseWriter, request *http.Request) {
	sort := request.FormValue("sort")
	if sort == "" {
		sort = "recently"
	}

	db, err := getMysqlConnection()
	if err != nil {
		http.Error(writer, fmt.Sprintf("cannot connect to database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT q.id, q.path, q.title, q.detail, q.asked_at, q.asked_by, u.name FROM questions q LEFT JOIN users u ON q.asked_by = u.username ORDER BY q.asked_at DESC")
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
		questions[i].Url = fmt.Sprintf("http://%s/questions/%d/%s", request.Host, question.Id, question.Path)
		calcPeriod(&questions[i])
	}

	tmpl := template.Must(template.ParseFiles("template/questions.html"))
	_ = tmpl.Execute(writer, QuestionsView{
		Questions:         questions,
		SortByRecently:    sort == "recently",
		SortByInteresting: sort == "interesting",
	})
}

func getMysqlConnection() (*sql.DB, error) {
	mysqlHost := os.Getenv(MysqlHost)
	mysqlUsername := os.Getenv(MysqlUsername)
	mysqlPassword := os.Getenv(MysqlPassword)
	return sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/iwantoask?parseTime=true", mysqlUsername, mysqlPassword, mysqlHost))
}

func calcPeriod(question *Question) {
	since := time.Since(question.AskedAt)
	if since.Seconds() < 60 {
		question.Since = fmt.Sprintf("%d seconds ago", int(since.Seconds()))
	} else if since.Minutes() < 60 {
		question.Since = fmt.Sprintf("%d minutes ago", int(since.Minutes()))
	} else if since.Hours() < 24 {
		question.Since = fmt.Sprintf("%d hours ago", int(since.Hours()))
	} else if since.Hours()/24 < 30 {
		question.Since = fmt.Sprintf("%d days ago", int(since.Hours()/24))
	} else if since.Hours()/24/30 < 30 {
		question.Since = fmt.Sprintf("%d months ago", int(since.Hours()/24/30))
	} else {
		question.Since = fmt.Sprintf("%d years ago", int(since.Hours()/24/30/12))
	}
}

func AskQuestion(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/ask.html"))
	_ = tmpl.Execute(writer, nil)
}

func SubmitQuestion(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Post Form: %s", request.PostFormValue("title"))
}
