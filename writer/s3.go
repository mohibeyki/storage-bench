package writer

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Writer struct {
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
}

func (w *S3Writer) WriteFile(path string, reader io.Reader) (err error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(w.Region),
		Credentials: credentials.NewStaticCredentials(w.AccessKey, w.SecretKey, "")},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	uploader := s3manager.NewUploader(sess)

	if _, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(w.Bucket),
		Key:    &path,
		Body:   reader,
	}); err != nil {
		fmt.Println(err)
		return err
	}

	return
}
