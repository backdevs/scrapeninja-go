package scrapeninja

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"time"
)

var (
	BaseUrl                    = "https://scrapeninja.p.rapidapi.com"
	ContentType                = "application/json"
	RapidAPIHost               = "scrapeninja.p.rapidapi.com"
	Timeout      time.Duration = 30
)

type ScrapeRequest struct {
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Data    string            `json:"body,omitempty"`
	Proxy   string            `json:"geo,omitempty"`
}

func (r *ScrapeRequest) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0)
	for header := range r.Headers {
		keys = append(keys, header)
	}

	sort.Strings(keys)

	headers := make([]string, 0)
	for _, header := range keys {
		headers = append(headers, fmt.Sprintf("%s: %s", header, r.Headers[header]))
	}

	type ScrapeRequestAlias ScrapeRequest
	return json.Marshal(&struct {
		*ScrapeRequestAlias
		Headers []string `json:"headers"`
	}{
		ScrapeRequestAlias: (*ScrapeRequestAlias)(r),
		Headers:            headers,
	})
}

type ScrapeResponse struct {
	Info struct {
		Version       string            `json:"version"`
		StatusCode    int               `json:"statusCode"`
		StatusMessage string            `json:"statusMessage"`
		Headers       map[string]string `json:"headers"`
	} `json:"info"`
	Body string `json:"body"`
}

type Client struct {
	httpClient *http.Client
	apiKey     string
}

func (c *Client) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	u, _ := url.Parse(BaseUrl)
	u.Path = path

	request, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", ContentType)
	request.Header.Set("X-RapidAPI-Host", RapidAPIHost)
	request.Header.Set("X-RapidAPI-Key", c.apiKey)

	return request, err
}

func (c *Client) Scrape(request *ScrapeRequest) (*ScrapeResponse, error) {
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	scrapeRequest, err := c.newRequest(
		http.MethodPost,
		"/scrape",
		bytes.NewBuffer(jsonRequest),
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(scrapeRequest)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var scrapeResponse ScrapeResponse

	err = json.Unmarshal(body, &scrapeResponse)
	if err != nil {
		return nil, err
	}

	return &scrapeResponse, nil
}

func New(apiKey string) Client {
	return Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: time.Second * Timeout},
	}
}
