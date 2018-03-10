package main

import (
	"testing"
	"net/http"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r* http.Request) {
		tpl.ExecuteTemplate(w, "index.gohtml", r)

		res, _ := http.Get("http://localhost:8000")
		if res.StatusCode != 200 {
			t.Errorf("Can't connect to localhost on port 8000")
		}

	})

}



