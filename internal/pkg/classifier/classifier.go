package classifier

import (
	"github.com/AsFal/shopify-application/internal/pkg/imgrepo"
)

type Img struct{}

type Classifier interface {
	Classify(imgrepo.ImgURI) (string, error)
}
