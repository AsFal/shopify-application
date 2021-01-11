package api

import (
	"log"
	"net/http"
	"os"

	"github.com/AsFal/shopify-application/internal/pkg/imgrepo"
	"github.com/AsFal/shopify-application/internal/pkg/imgrepo/amazon"
	"github.com/gin-gonic/gin"
)

func router() *gin.Engine {

	r := gin.Default()

	session, err := amazon.ConnectAws(
		os.Getenv("ACCESS_KEY_ID"),
		os.Getenv("ACCES_KEY"),
		os.Getenv("REGION"), // TODO: Change this to a constant
	)
	if err != nil {
		log.Println("The AmazonS3 credentials provided are missing or incorrect.")
		log.Println("The API will not support Upload functionality.")
	}

	// TODO: The Connection function should verify that the bucket is valid
	var client imgrepo.ImgRepoClient
	client = amazon.NewAmazonS3Client(session, os.Getenv("BUCKET"))

	r.POST("/", func(c *gin.Context) {
		fileHeader, err := c.FormFile("image")
		file, _ := fileHeader.Open() // TODO: Handle error
		if err != nil {
			c.String(http.StatusInternalServerError, "No image at 'image' form key")
		}
		client.Upload(file)
	})

	r.GET("/", func(c *gin.Context) {
		log.Println("Hit Get Index")
	})

	return r
}
