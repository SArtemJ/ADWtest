package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	srv         *httptest.Server
	slicePITest []PersonalIssue
)

func Test_init(t *testing.T) {

	if PersonalToken == "" {
		t.Errorf("Incorrect token type = %T value = %v", PersonalToken, PersonalToken)
	}

	if MyClient == nil {
		t.Errorf("Can't connect to github - no client")
	}
}

func Test_main(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer srv.Close()

	res, err := http.Get(srv.URL)
	if err != nil {
		log.Fatal(err)
	}

	value, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", value)
}

func Test_GetAllIssues(t *testing.T) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		req, err := http.Get("http://localhost:8000")
		if err != nil {
			t.Errorf("No request result status %v", req.Status)
		}

		req, err = http.PostForm("http://localhost:8000", url.Values{"statusI": {"testS"}, "labelI": {"testL"}, "repoI": {"testR"}})
		if err != nil {
			t.Errorf("Bad param to form %v %v %v %v", r.FormValue("statusI"), r.FormValue("labelI"), r.FormValue("repoI"), req.Status)
		}

	})

}

func Test_createPersonalIssues(t *testing.T) {

	opt := &github.IssueListByRepoOptions{
		State: Status,
		Labels: Label,
	}

	issues, _, err := MyClient.Issues.ListByRepo(context.Background(), "SArtemJ", Repo, opt)
	if err != nil {
		t.Errorf("No issues in Repo")
	}

	for _, i := range issues {
		pI := &PersonalIssue{
			ID:       *i.ID,
			Title:    *i.Title,
			Repo:     *i.RepositoryURL,
			Assignee: *i.Assignee.Login,
			Labels:   createSliceLabel(*i),
			Status:   *i.State,
		}
		slicePITest = append(slicePITest, *pI)
	}

	if len(slicePITest) == 0 {
		t.Errorf("No issues in Repo len of personal slice %v", len(slicePITest))
	}

}

func Test_createSliceLabel(t *testing.T) {
	var sL []string
	testLabel := []string{"test1", "test2"}
	for _, l := range testLabel {
		sL = append(sL, l)
	}

	if len(sL) == 0 {
		t.Errorf("Create slice of labels - error - zero len slice %v", len(sL))
	}
}
