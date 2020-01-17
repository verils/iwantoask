package app

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	haikunator "github.com/atrox/haikunatorgo/v2"
	"github.com/boltdb/bolt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const BucketQuestions = "questions"

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

type Question struct {
	Id      int       `json:"id"`
	Title   string    `json:"title"`
	Detail  string    `json:"detail"`
	Path    string    `json:"path"`
	AskedAt time.Time `json:"asked_at"`
	AskedBy User      `json:"asked_by"`
	Since   string    `json:"since,omitempty"`
}

type ListQuestionsView struct {
	BasePath          string
	Questions         []Question
	SortByRecently    bool
	SortByInteresting bool
	Pagination
}

type AskQuestionView struct {
	BasePath    string
	TitleError  string
	DetailError string
}

func (view *AskQuestionView) HasError() bool {
	return view.TitleError != "" || view.DetailError != ""
}

type Handler struct {
	db *bolt.DB
}

func NewHandler(db *bolt.DB) *Handler {
	return &Handler{db}
}

func (handler *Handler) ListQuestions(writer http.ResponseWriter, request *http.Request) {
	pageValue := FormValueOrDefault(request, "page", "1")
	page, err := strconv.Atoi(pageValue)
	if err != nil {
		http.Error(writer, "cannot parse param 'page' to int", http.StatusBadRequest)
		return
	}
	if page < 1 {
		http.Error(writer, "'page' cannot be < 1", http.StatusBadRequest)
		return
	}

	sizeValue := FormValueOrDefault(request, "size", "10")
	size, err := strconv.Atoi(sizeValue)
	if err != nil {
		http.Error(writer, "cannot parse param 'size' to int", http.StatusBadRequest)
		return
	}
	if size < 0 {
		http.Error(writer, "'size' cannot be < 0", http.StatusBadRequest)
		return
	}

	allQuestions := make([]Question, 0)

	err = handler.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketQuestions))
		cursor := bucket.Cursor()

		for k, v := cursor.Last(); k != nil; k, v = cursor.Prev() {
			var question Question
			err := json.Unmarshal(v, &question)
			if err != nil {
				return err
			}
			allQuestions = append(allQuestions, question)
		}
		return nil
	})

	if err != nil {
		message := fmt.Sprintf("list error: %s", err.Error())
		log.Printf("[ERROR]: %s", message)
		http.Error(writer, message, http.StatusInternalServerError)
		return
	}

	total := len(allQuestions)

	start := (page - 1) * size
	if start > total {
		start = total
	}

	end := page * size
	if end > total {
		end = total
	}

	questions := allQuestions[start:end]

	for i := range questions {
		calcPeriod(&questions[i])
	}

	pagination := NewPagination(page, size, total)
	pagination.Prepare()

	funcMap := template.FuncMap{"isActive": pagination.IsActive}
	tmpl := template.Must(template.New("list").Funcs(funcMap).ParseFiles("template/questions.html"))
	_ = tmpl.ExecuteTemplate(writer, "list", ListQuestionsView{
		BasePath:   BasePath,
		Questions:  questions,
		Pagination: *pagination,
	})
}

func (handler *Handler) ListQuestionsJson(writer http.ResponseWriter, _ *http.Request) {
	allQuestions := make([]Question, 0)

	err := handler.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketQuestions))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var question Question
			err := json.Unmarshal(v, &question)
			if err != nil {
				return err
			}
			allQuestions = append(allQuestions, question)
		}
		return nil
	})

	if err != nil {
		message := fmt.Sprintf("{\"message\": \"list error: %s\"}", err.Error())
		log.Printf("[ERROR]: %s", message)
		http.Error(writer, message, http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(writer).Encode(allQuestions)
}

func (handler *Handler) AskQuestion(writer http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/ask.html"))
	_ = tmpl.Execute(writer, AskQuestionView{BasePath: BasePath})
}

func (handler *Handler) SubmitQuestion(writer http.ResponseWriter, request *http.Request) {
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
		AskedAt: time.Now(),
		AskedBy: User{
			Username: username,
			Name:     strings.ReplaceAll(username, "_", " "),
		},
		Since: "",
	}

	err := handler.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketQuestions))
		sequence, _ := bucket.NextSequence()

		question.Id = int(sequence)

		bytes := itob(question.Id)
		data, err := json.Marshal(question)
		if err != nil {
			return err
		}

		return bucket.Put(bytes, data)
	})

	if err != nil {
		message := fmt.Sprintf("an error happened when accessing the resources: %s", err.Error())
		log.Printf("[ERROR]: %s", message)
		http.Error(writer, message, http.StatusInternalServerError)
		return
	}

	log.Printf("[DEBUG] added question: %s", question.Title)

	writer.Header().Set("Location", BasePathPrefix("/questions"))
	writer.WriteHeader(http.StatusFound)
}

func itob(id int) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(id))
	return bytes
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
