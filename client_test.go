package scrapeninja

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	apiKey := "test-api-key"
	timeout := time.Second * Timeout

	client := New(apiKey)

	if client.apiKey != apiKey {
		t.Errorf("got: %s, wanted: %s", client.apiKey, apiKey)
	}

	if client.httpClient.Timeout != timeout {
		t.Errorf("got: %q, wanted: %q", client.httpClient.Timeout, timeout)
	}
}

func TestScrapeRequest_MarshalJSON(t *testing.T) {
	tests := map[string]struct {
		url     string
		proxy   string
		headers map[string]string
		method  string
		want    string
	}{
		"url is empty": {
			"",
			"test-proxy",
			nil,
			"GET",
			"{\"url\":\"\",\"method\":\"GET\",\"geo\":\"test-proxy\",\"headers\":[]}",
		},
		"proxy is empty": {
			"https://test",
			"",
			nil,
			"GET",
			"{\"url\":\"https://test\",\"method\":\"GET\",\"headers\":[]}",
		},
		"method is empty": {
			"https://test",
			"test-proxy",
			nil,
			"",
			"{\"url\":\"https://test\",\"method\":\"\",\"geo\":\"test-proxy\",\"headers\":[]}",
		},
		"headers is nil": {
			"https://test",
			"test-proxy",
			nil,
			"GET",
			"{\"url\":\"https://test\",\"method\":\"GET\",\"geo\":\"test-proxy\",\"headers\":[]}",
		},
		"headers has one element": {
			"https://test",
			"test-proxy",
			map[string]string{"test-header": "test-value"},
			"GET",
			"{\"url\":\"https://test\",\"method\":\"GET\",\"geo\":\"test-proxy\",\"headers\":[\"test-header: test-value\"]}",
		},
		"headers has multiple elements": {
			"https://test",
			"test-proxy",
			map[string]string{
				"test-header":  "test-value",
				"test-header2": "test-value2",
				"test-header3": "test-value3",
			},
			"GET",
			"{\"url\":\"https://test\",\"method\":\"GET\",\"geo\":\"test-proxy\",\"headers\":[\"test-header: test-value\",\"test-header2: test-value2\",\"test-header3: test-value3\"]}",
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			request := ScrapeRequest{
				Url:     data.url,
				Proxy:   data.proxy,
				Headers: data.headers,
				Method:  data.method,
			}

			got, err := request.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			if string(got) != data.want {
				t.Errorf("got: %s, wanted: %s", string(got), data.want)
			}
		})
	}
}

func TestClient_Scrape(t *testing.T) {
	apiKey := "test-api-key"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(&ScrapeResponse{
			Info: struct {
				Version       string            `json:"version"`
				StatusCode    int               `json:"statusCode"`
				StatusMessage string            `json:"statusMessage"`
				Headers       map[string]string `json:"headers"`
			}{
				"test-version",
				200,
				"test-message",
				nil,
			},
			Body: "test-body",
		})
		if err != nil {
			t.Fatal(err)
		}

		t.Run("path is correct", func(t *testing.T) {
			got := r.URL.Path
			if got != "/scrape" {
				t.Errorf("got: %s, wanted: %s", got, "scrape")
			}
		})

		t.Run("X-RapidAPI-Key header is set correctly", func(t *testing.T) {
			got := r.Header.Get("X-RapidAPI-Key")
			if got != apiKey {
				t.Errorf("got: %s, wanted: %s", got, apiKey)
			}
		})

		t.Run("X-RapidAPI-Host header is set correctly", func(t *testing.T) {
			got := r.Header.Get("X-RapidAPI-Host")
			if got != RapidAPIHost {
				t.Errorf("got: %s, wanted: %s", got, RapidAPIHost)
			}
		})

		t.Run("Content-Type header is set correctly", func(t *testing.T) {
			got := r.Header.Get("Content-Type")
			if got != ContentType {
				t.Errorf("got: %s, wanted: %s", got, ContentType)
			}
		})
	}))
	defer srv.Close()

	// Overwrite base url with test server url
	BaseUrl = srv.URL

	client := New(apiKey)
	response, _ := client.Scrape(&ScrapeRequest{
		Url:     "",
		Proxy:   "",
		Headers: nil,
		Method:  "",
	})

	if response.Body != "test-body" {
		t.Errorf("Oh no!")
	}
}
