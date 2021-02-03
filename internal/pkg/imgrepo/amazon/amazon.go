package amazon

import (
	"mime/multipart"
	"fmt"

	"github.com/AsFal/shopify-application/internal/pkg/imgrepo"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/google/uuid"
)

type AmazonS3Client struct {
	session *session.Session
	bucket  string
}

func NewAmazonS3Client(region string, bucket string) (*AmazonS3Client, error) {
	s, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
		})

	return &AmazonS3Client{
		session: s,
		bucket:  bucket,
	}, err
}

func (c *AmazonS3Client) Upload(file multipart.File) (imgrepo.ImgURI, error) {
	s3Svc := s3.New(c.session)
	imageName := fmt.Sprintf("image_%s.jpg", uuid.New().String())

	// Create an uploader with S3 client and default options
	uploader := s3manager.NewUploaderWithClient(s3Svc)

	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(imageName),
		Body: file,
	})

	if err != nil {
		return "", err
	}

	return imgrepo.ImgURI(up.Location), err
}
