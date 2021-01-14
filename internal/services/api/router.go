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

		var tags []imgrepo.ImgURI

		fileHeader, err := c.FormFile("image")
		if err != nil {
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
		}

		tagsJson := c.QueryMap("tags")
		if tags != nil {
			buf := new(bytes.Buffer)
			if err := json.NewEncoder(buf).Encode(tags); err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return 
			}
		} 

		// text := c.QueryMap("text")
		// if text != nil {
		// 	tags, err = s.Tokenizer.process(text)
		// 	c.String(http.StatusInternalServerError, err.Error())
		// 	return 
		// }

		imgUris, err = s.SearchClient.SearchByTag(strings.Fields(tags))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return 
		}
		c.JSON(http.StatusOK, imgUris)
	})

	return r
}
