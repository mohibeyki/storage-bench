package writer

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Writer struct {
	Bucket  string
	Session *session.Session
}

func (w *S3Writer) WriteFile(path string, reader io.Reader) (err error) {
	uploader := s3manager.NewUploader(w.Session)

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
