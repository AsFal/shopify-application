package local

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"

	"github.com/AsFal/shopify-application/internal/pkg/imgrepo"
	"github.com/google/uuid"
)

type LocalImgRepo struct {
	path string
}

func NewLocalImgRepo(path string) *LocalImgRepo {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir)
	}

	return &LocalImgRepo{
		path: path,
	}
}

func (c *LocalImgRepo) Upload(file multipart.File) (imgrepo.ImgURL, error) {

	imageName := fmt.Sprintf("image_%s", uuid.New().String())
	imagePath := path.Join(c.path, imageName)

	imageFile, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	defer imageFile.Close()

	io.Copy(imageFile, file)

	return imgrepo.ImgURL(imagePath), nil
}
