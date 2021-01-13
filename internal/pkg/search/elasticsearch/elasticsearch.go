package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/AsFal/shopify-application/internal/pkg/imgrepo"
	"github.com/AsFal/shopify-application/internal/pkg/search"
	"github.com/elastic/go-elasticsearch/v7"
)

type ElasticsearchSearch struct {
	baseUrl *url.URL
	client  *http.Client
}

const indexName = "images"

func NewElasticsearchSearch(host string) *ElasticsearchSearch {
	return &ElasticsearchSearch{
		baseUrl: &url.URL{
			Scheme: "http",
			Host:   host,
		},
		client: &http.Client{},
	}
}

func (es *ElasticsearchSearch) SearchByTag(tags []string) ([]imgrepo.ImgURL, error) {
	c, _ := elasticsearch.NewDefaultClient()

	buf := new(bytes.Buffer)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "test",
			},
		},
	}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, _ := c.Search(
		c.Search.WithContext(context.Background()),
		c.Search.WithIndex("image"),
		c.Search.WithBody(buf),
		c.Search.WithPretty(),
	)

	defer res.Body.Close()
	var e map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(e); err != nil {
		return nil, err
	}

	return nil, nil
}

func (es *ElasticsearchSearch) IndexImgData(data *search.ImgData) error {
	cfg := elasticsearch.Config{
		Addresses: []string{
			es.baseUrl.String(),
		},
	  }
	c , err := elasticsearch.NewClient(cfg)

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return err
	}

	_, err = c.Index(
		indexName,
		buf,
		c.Index.WithContext(context.Background()),
	)
	return err
}
