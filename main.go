package main

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"html/template"
	"log"
	"net/http"
	"flag"
)

//структура задачи, только с нужными нам полями
type personalIssue struct {
	ID       int64
	Title    string
	Repo     string
	Assignee string
	Labels   []string
	Status   string
}

//

//глобальные переменные
var (
	PersonalToken = ""
	MyClient      = github.NewClient(nil)
	tpl           *template.Template
	Status        = "open"
	Label         = "bug"
	Repo = "ADWtest"
)

func init() {

	//параметр ключа можно задавать при запуске, если не указываем используется ключ по умолчанию
	tk := flag.String("token", "", "")
	flag.Parse()
	PersonalToken = *tk

	//используемые шаблоны веб страниц
	tpl = template.Must(template.ParseGlob("templates/*"))

	//создаем нового клиента github с персональным токеном
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: PersonalToken},
	)
	tokenClient := oauth2.NewClient(context.Background(), tokenService)
	MyClient = github.NewClient(tokenClient)
}


func main() {

	//запускаем сервис на 8000 порту
	//так как страница только одна сперва вызывается функция получения задач - базовая
	http.HandleFunc("/", GetAllIssues)
	http.ListenAndServe(":8000", nil)
}


func GetAllIssues(w http.ResponseWriter, r *http.Request) {

	//следим за тем, какие методы выполняются на форме
	switch r.Method {
	case "GET":
		//по умолчанию получаем все задачи с заданными параметрами по умолчниаю open / bug
		defaultForm(w, r)
	case "POST":
		//когда отправляем данные с формы запускаем функцию отбора задач по критериям
		Status = r.FormValue("statusI")
		Label = r.FormValue("labelI")
		Repo = r.FormValue("repoI")
		filteredForm(w, r)
	}
}

func defaultForm(w http.ResponseWriter, r *http.Request) {

	//получаем задачи и переводим в удобную нам структуру
	t := createPersonalIssues()
	//передаем структуру в шаблон для отображения
	tpl.ExecuteTemplate(w, "index.gohtml", t)
}

func filteredForm(w http.ResponseWriter, r *http.Request) {

	t := createPersonalIssues()
	//Полученные задачи на этапе старта сервиса проверяем по условиям
	result := checkIssues(t)
	tpl.ExecuteTemplate(w, "index.gohtml", result)
}


func createPersonalIssues() []personalIssue {

	//получаем все задачи
	issues, _, err := MyClient.Issues.ListByRepo(context.Background(), "SArtemJ", Repo, nil)
	if err != nil {
		log.Panic("No issues in Repo")
	}

	var slicePI []personalIssue
	for _, i := range issues {
		pI := &personalIssue{
			ID:       *i.ID,
			Title:    *i.Title,
			Repo:     *i.RepositoryURL,
			Assignee: *i.Assignee.Login,
			Labels:   createSliceLabel(*i),
			Status:   *i.State,
		}
		slicePI = append(slicePI, *pI)
	}

	return slicePI
}

//так как у задачи может быть несколько лэйблов создаем удобную структуру
func createSliceLabel(in github.Issue) []string {
	var sL []string
	for _, l := range in.Labels {
		sL = append(sL, *l.Name)
	}
	return sL
}

//проверям задачи по статусу
func checkIssues(in []personalIssue) []personalIssue {

	var t []personalIssue
	var t2 []personalIssue
	for _, v := range in {
		if v.Status == Status {
			if checkLabels(v.Labels) {
				t = append(t, v)
			} else {
				t2 = append(t2, v)
			}
		}
	}

	if len(t) == 0 {
		return t2
	} else {
		return t
	}
}

//проверям задачу по лэйблу
//елси находится хотя бы один совпадающий - задача подходит для вывода
func checkLabels(in []string) bool {
	var t = false
	for _, l := range in {
		log.Println(Label)
		if l == (Label) {
			t = true
			break
		}
	}
	return t
}