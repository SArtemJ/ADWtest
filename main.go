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

type personalIssue struct {
	ID       int64
	Title    string
	Repo     string
	Assignee string
	Labels   []string
	Status   string
}

//"e36cdea8137a975e1a67a0084b80eac079146fc8"

var (
	PersonalToken = ""
	MyClient      = github.NewClient(nil)
	tpl           *template.Template
	Status        = "open"
	Label         = "bug"
	T []personalIssue
)

func init() {

	tk := flag.String("token", "e36cdea8137a975e1a67a0084b80eac079146fc8", "")
	flag.Parse()
	PersonalToken = *tk

	tpl = template.Must(template.ParseGlob("templates/*"))

	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: PersonalToken},
	)
	tokenClient := oauth2.NewClient(context.Background(), tokenService)
	MyClient = github.NewClient(tokenClient)
}

func GetAllIssues(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		defaultForm(w, r)
	case "POST":
		Status = r.FormValue("statusI")
		Label = r.FormValue("labelI")
		filteredForm(w, r)
	}
}

func defaultForm(w http.ResponseWriter, r *http.Request) {
	issues, _, err := MyClient.Issues.ListByRepo(context.Background(), "SArtemJ", "ADWtest", nil)
	if err != nil {
		log.Panic("No issues in Repo")
	}

	T = createPersonalIssues(issues)
	tpl.ExecuteTemplate(w, "index.gohtml", T)
}

func filteredForm(w http.ResponseWriter, r *http.Request) {
	//issues, _, err := MyClient.Issues.ListByRepo(context.Background(), "SArtemJ", "ADWtest", nil)
	//if err != nil {
	//	//log.Panic("No issues in Repo")
	//}

	result := checkIssues(T)
	tpl.ExecuteTemplate(w, "index.gohtml", result)
}

func main() {

	http.HandleFunc("/", GetAllIssues)
	http.ListenAndServe(":8000", nil)
}


func createSliceLabel(in github.Issue) []string {

	var sL []string
	for _, l := range in.Labels {
		sL = append(sL, *l.Name)
	}

	return sL

}


func createPersonalIssues(in []*github.Issue) []personalIssue {

	var slicePI []personalIssue
	for _, i := range in {
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