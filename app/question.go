package app

import (
	"database/sql"
	"fmt"
	"github.com/atrox/haikunatorgo/v2"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
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
	Since   string
}

type ListQuestionsView struct {
	BasePath          string
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
		message := fmt.Sprintf("cannot connect to database: %s", err.Error())
		log.Printf("[ERROR]: %s", message)
		http.Error(writer, message, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT q.id, q.path, q.title, q.detail, q.asked_at, q.asked_by, u.name FROM questions q LEFT JOIN users u ON q.asked_by = u.username ORDER BY q.asked_at DESC")
	if err != nil {
		message := fmt.Sprintf("cannot fetch questions: %s", err.Error())
		log.Printf("[ERROR]: %s", message)
		http.Error(writer, message, http.StatusInternalServerError)
		return
	}

	questions := make([]Question, 0)
	for rows.Next() {
		var question Question
		var userModel UserModel

		err = rows.Scan(&question.Id, &question.Path, &question.Title, &question.Detail, &question.AskedAt, &userModel.Username, &userModel.Name)
		if err != nil {
			message := fmt.Sprintf("scan error: %s", err.Error())
			log.Printf("[ERROR]: %s", message)
			http.Error(writer, message, http.StatusInternalServerError)
			return
		}

		question.AskedBy = User{
			Username: userModel.GetUsernameOrDefault("anonymous"),
			Name:     userModel.GetNameOrDefault("anonymous"),
		}
		questions = append(questions, question)
	}

	for i, question := range questions {
		questions[i].Url = getUrl(request, question)
		calcPeriod(&questions[i])
	}

	tmpl := template.Must(template.ParseFiles("template/questions.html"))
	_ = tmpl.Execute(writer, ListQuestionsView{
		BasePath:          BasePath,
		Questions:         questions,
		SortByRecently:    sort == "recently",
		SortByInteresting: sort == "interesting",
	})
}

func getMysqlConnection() (*sql.DB, error) {
	mysqlHost := os.Getenv(EnvMysqlHost)
	mysqlUsername := os.Getenv(EnvMysqlUsername)
	mysqlPassword := os.Getenv(EnvMysqlPassword)
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

type AskQuestionView struct {
	BasePath    string
	TitleError  string
	DetailError string
}

func (view *AskQuestionView) HasError() bool {
	return view.TitleError != "" || view.DetailError != ""
}

func AskQuestion(writer http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/ask.html"))
	_ = tmpl.Execute(writer, AskQuestionView{BasePath: BasePath})
}

func SubmitQuestion(writer http.ResponseWriter, request *http.Request) {
	view := AskQuestionView{BasePath: BasePath}

	title := request.PostFormValue("title")
	if title == "" {
		view.TitleError = "请填写标题"
	}
	detail := request.PostFormValue("detail")
	if detail == "" {
		view.DetailError = "请填写描述"
	}

	if view.HasError() {
		askTemplate := template.Must(template.ParseFiles("template/ask.html"))
		_ = askTemplate.Execute(writer, view)
		return
	}

	haikunate := haikunator.New()
	haikunate.Delimiter = "_"
	haikunate.TokenLength = 0
	username := haikunate.Haikunate()

	path := strings.ReplaceAll(strings.ToLower(title), " ", "+")

	question := Question{
		Title:   title,
		Detail:  detail,
		Path:    path,
		Url:     "",
		AskedAt: time.Now(),
		AskedBy: User{
			Username: username,
			Name:     strings.ReplaceAll(username, "_", " "),
		},
		Since: "",
	}

	db, err := getMysqlConnection()
	if err != nil {
		message := fmt.Sprintf("cannot connect to database: %s", err.Error())
		log.Printf("[ERROR]: %s", message)
		http.Error(writer, message, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO questions (path, title, detail, asked_at, asked_by) VALUES (?, ?, ?, ?, ?)",
		question.Path, question.Title, question.Detail, question.AskedAt, question.AskedBy.Username)
	if err != nil {
		message := fmt.Sprintf("an error happened when accessing the resources: %s", err.Error())
		log.Printf("[ERROR]: %s", message)
		http.Error(writer, message, http.StatusInternalServerError)
		return
	}
	lastInsertId, _ := result.LastInsertId()

	question.Id = int(lastInsertId)

	question.Url = getUrl(request, question)

	writer.Header().Set("Location", BasePathPrefix("/questions"))
	writer.WriteHeader(http.StatusFound)
}

func getUrl(request *http.Request, question Question) string {
	schema := "http"
	if request.TLS != nil {
		schema = "https"
	}
	return fmt.Sprintf("%s://%s/%s/questions/%d/%s", schema, request.Host, BasePath, question.Id, question.Path)
}
