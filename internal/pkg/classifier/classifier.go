package classifier

import (
	"mime/multipart"

	"github.com/AsFal/shopify-application/internal/pkg/search"
)

type Img struct{}

type Classifier interface {
	classify(multipart.File) search.ImgData
}
