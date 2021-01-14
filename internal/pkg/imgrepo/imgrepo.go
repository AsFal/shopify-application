package imgrepo

import "mime/multipart"

type ImgURI string

type ImgRepoClient interface {
	Upload(multipart.File) (ImgURI, error)
}
