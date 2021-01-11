package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/elastic/go-elasticsearch/v7"
)

type ElasticSearchTokenizer struct {
	baseUrl *url.URL
	client  *http.Client
}

func NewTokenizer(host string) *ElasticSearchTokenizer {
	return &ElasticSearchTokenizer{
		baseUrl: &url.URL{
			Scheme: "https",
			Host:   host,
		},
		client: &http.Client{},
	}
}

func (t *ElasticSearchTokenizer) Process(text string) ([]string, error) {
	client, _ := elasticsearch.NewDefaultClient()

	var buf bytes.Buffer
	body := map[string]interface{}{
		"analyzer": "standard",
		"text":     text,
	}

	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, _ := client.Indices.Analyze(
		client.Indices.Analyze.WithContext(context.Background()),
		client.Indices.Analyze.WithBody(&buf),
		client.Indices.Analyze.WithIndex(""),
	)

	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	return nil, nil
}
