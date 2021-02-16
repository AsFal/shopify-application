package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AsFal/shopify-application/internal/pkg/search"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Service) router() *gin.Engine {

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/image", func(c *gin.Context) {
		fileHeader, err := c.FormFile("image")
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, "%s", c.Errors.String())
			return
		}
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
		c.String(http.StatusCreated, string(uri))
	})

	r.POST("/_search/_image", func(c *gin.Context) {
		var tags []string

		fileHeader, err := c.FormFile("image")

		// Search With Similar Image
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
			tags = strings.Fields(string(tagsString))
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}

		// Input Gets Converted To Tags, Which Are Then Used to Query the Search Client
		imgUris, err := s.SearchClient.SearchByTag(tags)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, imgUris)
	})

	r.GET("/_search", func(c *gin.Context) {
		var tags []string

		// Search With Tags
		tagsJson := c.Query("tags")
		if tagsJson != "" {
			err := json.Unmarshal([]byte(tagsJson), &tags)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}

		// Search With Full Text
		text := c.Query("text")
		if text != "" {
			tags = strings.Fields(text)
		}

		// Input Gets Converted To Tags, Which Are Then Used to Query the Search Client
		imgUris, err := s.SearchClient.SearchByTag(tags)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, imgUris)
	})

	return r
}
