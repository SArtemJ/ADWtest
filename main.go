package main

import (
	"context"
	"flag"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//структура задачи, только с нужными нам полями
type PersonalIssue struct {
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
	//token
	PersonalToken = ""
	//client github
	MyClient = github.NewClient(nil)
	//html templates
	tpl *template.Template
	// status for issues
	Status = "all"
	//labels issues
	Label = []string{}
	//repo issues
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

	Label = nil
}

func main() {

	//запускаем сервис на 8000 порту
	//так как страница только одна сперва вызывается функция получения задач - базовая
	http.HandleFunc("/", GetAllIssues)
	http.ListenAndServe(":8000", nil)
}

func GetAllIssues(w http.ResponseWriter, r *http.Request) {

	//следим за тем, какие методы выполняются
	switch r.Method {
	case "GET":
		//по умолчанию получаем все задачи с заданными параметрами по умолчниаю all / all
		defaultForm(w, r)
	case "POST":
		//когда отправляем данные с формы запускаем функцию отбора задач по критериям
		if r.FormValue("statusI") != "" {
			Status = r.FormValue("statusI")
		} else {Status = "all"}
		if r.FormValue("labelI") != "" {
			Label = nil
			//так как лэйблов может быть несколько, сначала обнуляем слайс
			// и заполняем значениями из формы
			Label = strings.Split(r.FormValue("labelI"), ",")
		} else {Label = nil}
		if r.FormValue("repoI") != "" {
			Repo = r.FormValue("repoI")
		}
		defaultForm(w, r)
	}
}

func defaultForm(w http.ResponseWriter, r *http.Request) {

	//получаем задачи и переводим в удобную нам структуру
	t := createPersonalIssues()
	//передаем структуру в шаблон для отображения
	tpl.ExecuteTemplate(w, "index.gohtml", t)
}

func createPersonalIssues() []PersonalIssue {

	//опциональный параметр библиотеки позволяет отобрать задачи с заданным статусом и лйэблами
	opt := &github.IssueListByRepoOptions{
		State:  Status,
		Labels: Label,
	}

	//получаем все задачи
	issues, _, err := MyClient.Issues.ListByRepo(context.Background(), "SArtemJ", Repo, opt)
	if err != nil {
		log.Panic("No issues in Repo")
	}

	//для отображения перепишем в удобный слайс
	var slicePI []PersonalIssue
	for _, i := range issues {
		pI := &PersonalIssue{
			ID:       *i.ID,
			Title:    *i.Title,
			Repo:     *i.RepositoryURL,
			Assignee: *i.Assignee.Login,
			//так как объет label с несколькими полями, нас интересует только имя лэйбда
			Labels: createSliceLabel(*i),
			Status: *i.State,
		}
		slicePI = append(slicePI, *pI)
	}

	return slicePI
}

//так как у задачи может быть несколько лэйблов создаем удобную структуру
// получаем только имена лэйблов
func createSliceLabel(in github.Issue) []string {
	var sL []string
	for _, l := range in.Labels {
		sL = append(sL, *l.Name)
	}
	return sL
}

