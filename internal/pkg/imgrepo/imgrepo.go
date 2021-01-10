package imgrepo

import "mime/multipart"

type ImgURL string

type ImgRepoClient interface {
	upload(multipart.FileHeader) ImgURL
}
