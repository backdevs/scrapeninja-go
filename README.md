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
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	ninja := scrapeninja.New("Your-API-Key")
	response, err := ninja.Scrape(&scrapeninja.ScrapeRequest{
		Url:   "https://jsonplaceholder.typicode.com/posts/1",
		Proxy: "us",
		Headers: map[string]string{
			"Accept": "application/json",
		},
		Method: http.MethodGet,
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