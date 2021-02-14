package classifier

import (
	"github.com/AsFal/shopify-application/internal/pkg/imgrepo"
	"github.com/AsFal/shopify-application/internal/pkg/search"
)

type Img struct{}

type Classifier interface {
	Classify(imgrepo.ImgURI) (search.Tags, error)
}
