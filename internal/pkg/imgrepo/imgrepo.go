package imgrepo

import "mime/multipart"

type ImgURL string

type ImgRepoClient interface {
	Upload(multipart.File) (ImgURL, error)
}
