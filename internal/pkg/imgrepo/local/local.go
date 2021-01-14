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
	dockerHostVolumePath string
}

func NewLocalImgRepo(path string, dockerHostVolumePath string) *LocalImgRepo {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir)
	}

	return &LocalImgRepo{
		path: path,
		dockerHostVolumePath: dockerHostVolumePath,
	}
}

func (c *LocalImgRepo) Upload(file multipart.File) (imgrepo.ImgURI, error) {
	imageName := fmt.Sprintf("image_%s.jpg", uuid.New().String())
	imagePath := path.Join(c.path, imageName)

	imageFile, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	defer imageFile.Close()

	io.Copy(imageFile, file)

	if c.dockerHostVolumePath == "" {
		return imgrepo.ImgURI(imagePath), nil
	} else {
		return imgrepo.ImgURI(path.Join(c.dockerHostVolumePath, imageName)), nil
	}
}
