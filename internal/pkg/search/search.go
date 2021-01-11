package search

import "github.com/AsFal/shopify-application/internal/pkg/imgrepo"

type ImgData struct {
	Url  imgrepo.ImgURL
	Tags string
}

type SearchClient interface {
	SearchByTag([]string) ([]imgrepo.ImgURL, error)
	IndexImgData(*ImgData) error
}
