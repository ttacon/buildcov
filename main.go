package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/ttacon/pretty"
)

var matcher = regexp.MustCompile("coverage: ([\\d]+\\.[\\d]+)% of statements")

func main() {
	http.HandleFunc("/travisci", handleBuild)
	http.ListenAndServe(":18009", nil)
}

func handleBuild(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	byt, err := ioutil.ReadAll(r.Body)
	fmt.Println("body was: ", string(byt))
	fmt.Println("==========")
	fmt.Printf("req: %#v\n", r)
	fmt.Println("==========")
	if err != nil {
		w.WriteHeader(http.StatusOK) // tell travis everything is okay
		return
	}

	var data TravisCIWebHookNotification
	err = json.Unmarshal(byt, &data)
	fmt.Println("err: ", err)
	pretty.Println(data)
	w.WriteHeader(http.StatusOK)
}

type TravisCIWebHookNotification struct {
	Payload *TravisCIPayload `json:"payload"`
}

type TravisCIPayload struct {
	ID             int        `json:""`
	Number         string     `json:"number"`
	Status         *string    `json:"status"`
	StartedAt      *time.Time `json:"started_at"`
	FinishedAt     *time.Time `json:"finished_at"`
	StatusMessage  string     `json:"status_message"`
	Commit         string     `json:"commit"`
	Branch         string     `json:"master"`
	Message        string     `json:"message"`
	CompareUrl     string     `json:"compare_url"`
	CommittedAt    *time.Time `json:"committed_at"`
	CommitterName  string     `json:"committer_name"`
	CommitterEmail string     `json:"committer_email"`
	AuthorName     string     `json:"author_name"`
	AuthorEmail    string     `json:"author_email"`
	Type           string     `json:"type"`
	BuildUrl       string     `json:"build_url"`
	// TODO(ttacon): there's a lot more here but i don't need it right now
}

/*
{
  "repository": {
    "id": 1,
    "name": "minimal",
    "owner_name": "svenfuchs",
    "url": "http://github.com/svenfuchs/minimal"
   },
  "config": {
    "notifications": {
      "webhooks": ["http://evome.fr/notifications", "http://example.com/"]
    }
  },
  "matrix": [
    {
      "id": 2,
      "repository_id": 1,
      "number": "1.1",
      "state": "created",
      "started_at": null,
      "finished_at": null,
      "config": {
        "notifications": {
          "webhooks": ["http://evome.fr/notifications", "http://example.com/"]
        }
      },
      "status": null,
      "log": "",
      "result": null,
      "parent_id": 1,
      "commit": "62aae5f70ceee39123ef",
      "branch": "master",
      "message": "the commit message",
      "committed_at": "2011-11-11T11: 11: 11Z",
      "committer_name": "Sven Fuchs",
      "committer_email": "svenfuchs@artweb-design.de",
      "author_name": "Sven Fuchs",
      "author_email": "svenfuchs@artweb-design.de",
      "compare_url": "https://github.com/svenfuchs/minimal/compare/master...develop"
    }
  ]
}
*/
