package main

import (
	"context"
	"testing"
	"log"
)

func test_access(t *testing.T) {

	PersonalToken = ""
	access()
	if MyClient == nil {
		t.Errorf("Something wrong with connect")
	}

}

func test_findIssues(t *testing.T) {

	PersonalToken = "00afb866ab29487921c547e1e7e67df50c3274de"
	access()

	issues, _, err := MyClient.Issues.ListByRepo(context.Background(), "SArtemJ", "ADWtest", nil)
	if err != nil {
		t.Errorf("Issues.List returned error: %v", err)
	}
	log.Println(len(issues))

}
