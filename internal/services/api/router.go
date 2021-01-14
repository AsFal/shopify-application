package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"log"

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
			return
		}
		uri, err := s.ImgRepoClient.Upload(file)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		tags, err := s.Classifier.Classify(uri)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		imgData := &search.ImgData{
			URI:  uri,
			Tags: tags,
		}
		err = s.SearchClient.IndexImgData(imgData)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.String(http.StatusOK, "OK")
	})

	r.POST("/_search", func(c *gin.Context) {

		var tags []string

		fileHeader, err := c.FormFile("image")

		if err == nil {
			file, _ := fileHeader.Open() // TODO: Handle error
			if err != nil {
				c.String(http.StatusInternalServerError, "No image at 'image' form key")
				return
			}
			uri, err := s.ImgRepoClient.Upload(file)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			tagsString, err := s.Classifier.Classify(uri)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			log.Println(tagsString)
			tags = strings.Fields(tagsString)
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}

		tagsJson := c.Query("tags")
		if tagsJson != "" {
			err := json.Unmarshal([]byte(tagsJson), &tags)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}

		text := c.Query("text")
		if text != "" {
			tags = strings.Fields(text)
		}

		imgUris, err := s.SearchClient.SearchByTag(tags)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, imgUris)
	})

	return r
}
