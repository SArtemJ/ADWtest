package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
)

//"00afb866ab29487921c547e1e7e67df50c3274de"

var (
	PersonalToken = ""
	Status        = ""
	Label         = ""
	MyClient      = github.NewClient(nil)
	Own = ""
	RP = ""
)

func init() {
	// ключи для запуска программы, необязательны
	// токен пользователя
	tk := flag.String("token", "00afb866ab29487921c547e1e7e67df50c3274de", "")
	// лэйбл для поиска задач
	lbl := flag.String("lbl", "bug", "")
	// статус задач
	sts := flag.String("sts", "open", "")
	// владелец
	own := flag.String("own", "SArtemJ", "")
	// repo
	rp := flag.String("rp", "ADWtest", "")

	flag.Parse()
	PersonalToken = *tk
	Status = *sts
	Label = *lbl
	Own = *own
	RP = *rp
}

func main() {
	access()
	findIssues()
}

func access() {
	//Создаем нового клинета github c персональным токеном
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: PersonalToken},
	)
	tokenClient := oauth2.NewClient(context.Background(), tokenService)
	MyClient = github.NewClient(tokenClient)
}

func findIssues() {

	//получаем все задачи из заданного репозитория
	//владелец и репозиторий задан по умолчанию, можно также поменять при запуске

	issues, _, err := MyClient.Issues.ListByRepo(context.Background(), Own, RP, nil)
	if err != nil {
		log.Panic("No issues in Repo")
	}
	for _, v := range issues {

		var stringLabels = ""
		for i := 0; i < len(v.Labels); i++ {

			// проверям чтобы статус задачи и лэйбл соответствовали ключам при запуске программы иначе ключи по умолчанию
			if *v.State == Status && *v.Labels[i].Name == Label {
				for _, k := range v.Labels {
					stringLabels = stringLabels + *k.Name
				}
				fmt.Printf("ID %v Title %v Repo %v Assignee %v Owner _ Labels %v Status %v \n", *v.ID, *v.Title, *v.RepositoryURL, *v.Assignee.Login, stringLabels, *v.State)
			}
		}
	}
}
