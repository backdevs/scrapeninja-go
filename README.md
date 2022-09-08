# [ScrapeNinja](https://scrapeninja.net) Web scraper Go API Client

A slim Go client for interacting with the ScrapeNinja API

## More info about [ScrapeNinja](https://scrapeninja.net)
* Official website - https://scrapeninja.net
* Rapid API documentation - https://rapidapi.com/restyler/api/scrapeninja

# Installation
```shell
go get -u github.com/backdevs/scrapeninja-go
```

# Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/backdevs/scrapeninja-go"
	"log"
	"net/http"
	"strings"
)

type Post struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	UserID    int      `json:"userId"`
	Tags      []string `json:"tags"`
	Reactions int      `json:"reactions"`
}

func main() {
	ninja := scrapeninja.New("Your-API-Key")

	postJson, err := json.Marshal(&Post{
		ID:     "1",
		UserID: 12,
		Title:  "ScrapeNinja Web scraper Go API Client",
		Body:   "Lorem ipsum etc.",
	})
	if err != nil {
		log.Fatal(err)
	}

	response, err := ninja.Scrape(&scrapeninja.ScrapeRequest{
		Url:    "https://dummyjson.com/posts/1",
		Method: http.MethodPut,
		Proxy:  scrapeninja.ProxyEU,
		Data:   string(postJson),
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	var post Post
	err = json.NewDecoder(strings.NewReader(response.Body)).Decode(&post)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(post.Title)
}

```