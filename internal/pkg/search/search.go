package search

import "github.com/AsFal/shopify-application/internal/pkg/imgrepo"

type ImgData struct {
	url  imgrepo.ImgURL
	tags []string
}

type SearchClient interface {
	searchByTag([]string) []imgrepo.ImgURL
}
