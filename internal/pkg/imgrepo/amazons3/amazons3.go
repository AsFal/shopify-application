package amazons3

import (
	"mime/multipart"

	"github.com/AsFal/shopify-application/internal/pkg/imgrepo"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AmazonS3Client struct{}

func NewAmazonS3Client() *AmazonS3Client {
	return &AmazonS3Client{}
}

func (*AmazonS3Client) upload(file multipart.File) imgrepo.ImgURL {
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)

	upParams := &s3manager.UploadInput{
		Bucket: "",
		Key:    "",
		Body:   file,
	}
	result, _ := uploader.Upload(upParams)
	return result
}
