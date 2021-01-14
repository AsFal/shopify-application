package search

import "github.com/AsFal/shopify-application/internal/pkg/imgrepo"

type ImgData struct {
	URI  imgrepo.ImgURI `json:"uri"`
	Tags string `json:"tags"`
}

type SearchClient interface {
	SearchByTag([]string) ([]imgrepo.ImgURI, error)
	IndexImgData(*ImgData) error
}
