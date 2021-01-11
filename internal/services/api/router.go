package api

import (
	"log"
	"net/http"

	"github.com/AsFal/shopify-application/internal/pkg/search"
	"github.com/gin-gonic/gin"
)

func (s *Service) router() *gin.Engine {

	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		fileHeader, err := c.FormFile("image")
		file, _ := fileHeader.Open() // TODO: Handle error
		if err != nil {
			c.String(http.StatusInternalServerError, "No image at 'image' form key")
		}
		url, err := s.ImgRepoClient.Upload(file)
		tags, err := s.Classifier.Classify(file)
		imgData := &search.ImgData{
			Url:  url,
			Tags: tags,
		}
		s.SearchClient.IndexImgData(imgData)
	})

	r.GET("/", func(c *gin.Context) {
		log.Println("Hit Get Index")
	})

	return r
}
