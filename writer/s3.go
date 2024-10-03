package writer

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Writer struct {
	Bucket string
	Client *s3.Client
}

func (w *S3Writer) WriteFile(path string, reader io.Reader) (err error) {
	uploader := manager.NewUploader(w.Client, func(u *manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024
	})

	if _, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(w.Bucket),
		Key:    aws.String(path),
		Body:   reader,
	}); err != nil {
		fmt.Println("MOHI!")
		fmt.Println(err)
		return err
	}

	return
}
