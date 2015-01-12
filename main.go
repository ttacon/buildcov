package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ttacon/pretty"
	"github.com/ttacon/travisci"
)

var matcher = regexp.MustCompile("coverage: ([\\d]+\\.[\\d]+)% of statements")
var travisToken = flag.String("tt", "", "travis CI token to use (only for testing)")

func main() {
	flag.Parse()

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

	clean, err := url.QueryUnescape(string(byt))
	if err != nil {
		w.WriteHeader(http.StatusOK) // tell travis everything is okay
		return
	}

	fmt.Println(strings.TrimPrefix(clean, "payload="))

	var data TravisCIPayload
	err = json.Unmarshal([]byte(strings.TrimPrefix(clean, "payload=")), &data)
	fmt.Println("err: ", err)
	pretty.Println(data)
	w.WriteHeader(http.StatusOK)

	go retrieveCoverageInfo(data.ID)
}

func retrieveCoverageInfo(id int) {
	fmt.Println("=====")

	c := travisci.NewClientFromTravis(*travisToken)
	build, err := c.GetBuildByID(strconv.Itoa(id))
	if err != nil {
		fmt.Println("failed to retrieve build:", id, ", err:", err)
		return
	}

	if len(build.JobIDs) == 0 {
		fmt.Println("no jobs were found for build:", id)
		return
	}

	// we only really care about one build right now
	// so grab the last one
	jID := build.JobIDs[len(build.JobIDs)-1]
	l, err := c.ArchivedLogByJob(jID)
	if err != nil {
		fmt.Println("failed to retreive log:", err)
		return
	}

	ms := matcher.FindStringSubmatch(string(l))
	fmt.Println(string(l[len(l)-200:]))
	fmt.Println(ms)
	fmt.Println("coverage: ", ms[1])
}

type TravisCIPayload struct {
	ID             int        `json:"id"`
	Number         string     `json:"number"`
	Status         *int       `json:"status"`
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
