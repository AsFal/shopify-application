package app

import "github.com/gin-gonic/gin"

func Router() *gin.Engine {
	r := gin.Default()
	r.POST("", func(c *gin.Context) {

	})

	r.GET("", func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			// Upload The Image
			// There is an image
		}
	})

	return r
}
