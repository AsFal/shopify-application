package classifier

import (
	"mime/multipart"
)

type Img struct{}

type Classifier interface {
	Classify(multipart.File) (string, error)
}
