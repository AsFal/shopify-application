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
		log.Println("line")
		fileHeader, err := c.FormFile("image")
		file, _ := fileHeader.Open() // TODO: Handle error
		if err != nil {
			c.String(http.StatusInternalServerError, "No image at 'image' form key")
		}
		url, err := s.ImgRepoClient.Upload(file)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		log.Println(url)
		tags, err := s.Classifier.Classify(file)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		log.Println(tags)
		imgData := &search.ImgData{
			Url:  url,
			Tags: tags,
		}
		err = s.SearchClient.IndexImgData(imgData)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.String(http.StatusOK, "OK")
	})

	r.GET("/", func(c *gin.Context) {
		log.Println("Hit Get Index")
	})

	return r
}
