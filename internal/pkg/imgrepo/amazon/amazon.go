package amazon

import (
	"mime/multipart"

	"github.com/AsFal/shopify-application/internal/pkg/imgrepo"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AmazonS3Client struct {
	session *session.Session
	bucket  string
}

func NewAmazonS3Client(session *session.Session, bucket string) *AmazonS3Client {
	return &AmazonS3Client{
		session: session,
		bucket:  bucket,
	}
}

func ConnectAws(
	accessKeyID string, accessKey string, region string,
) (*session.Session, error) {
	return session.NewSession(
		&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				accessKey,
				"", // a token will be created when the session it's used.
			),
		})
}

func (c *AmazonS3Client) Upload(file multipart.File) (imgrepo.ImgURI, error) {
	s3Svc := s3.New(c.session)

	// Create an uploader with S3 client and default options
	uploader := s3manager.NewUploaderWithClient(s3Svc)

	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(c.bucket),
		ACL:    aws.String("public-read"),
		// Key:    aws.String(filename),
		Body: file,
	})
	return imgrepo.ImgURI(up.Location), err
}
