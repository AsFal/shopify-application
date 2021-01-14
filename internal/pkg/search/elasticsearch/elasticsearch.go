package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

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

type SearchResponse struct {
	Hits HitWrapper `json:"hits"`
}

type HitWrapper struct {
	Hits []Hit `json:"hits"`
}

type Hit struct {
	Source search.ImgData `json:"_source"`
}

func (es *ElasticsearchSearch) SearchByTag(tags []string) ([]imgrepo.ImgURI, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			es.baseUrl.String(),
		},
	  }
	c , err := elasticsearch.NewClient(cfg)

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"tags": map[string]interface{}{
					"query": strings.Join(tags, " "),
					"minimum_should_match": 1,
				},
			},
		},
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := c.Search(
		c.Search.WithContext(context.Background()),
		c.Search.WithIndex("images"),
		c.Search.WithBody(buf),
		c.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	searchResponse := new(SearchResponse)
	if err := json.NewDecoder(res.Body).Decode(searchResponse); err != nil {
		return nil, err
	}


	imgUris := make([]imgrepo.ImgURI, 0)
	for _, hit := range searchResponse.Hits.Hits {
		imgUris = append(imgUris, hit.Source.URI)
	}

	return imgUris, nil
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
