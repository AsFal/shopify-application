package search

import (
	"strings"

	"github.com/AsFal/shopify-application/internal/pkg/imgrepo"
)

type ImgData struct {
	URI  imgrepo.ImgURI `json:"uri"`
	Tags Tags           `json:"tags"`
}

type Tags string

func (tags Tags) Contains(x string) bool {
	elements := strings.Split(string(tags), " ")
	for _, element := range elements {
		if element == x {
			return true
		}
	}
	return false
}

type SearchClient interface {
	SearchByTag([]string) ([]imgrepo.ImgURI, error)
	IndexImgData(*ImgData) error
}
