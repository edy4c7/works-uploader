package infrastructures

import (
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type StorageClientImpl struct {
	uploader   *s3manager.Uploader
	bucketName string
}

func NewStorageClientImpl() *StorageClientImpl {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return &StorageClientImpl{
		uploader:   s3manager.NewUploader(sess),
		bucketName: os.Getenv("S3_BUCKET"),
	}
}

func (r *StorageClientImpl) Upload(fileName string, fh *multipart.FileHeader) (string, error) {
	body, err := fh.Open()
	if err != nil {
		return "", err
	}

	uo, err := r.uploader.Upload(&s3manager.UploadInput{
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(fileName),
		Body:   body,
	})

	if uo != nil {
		return uo.Location, err
	}

	return "", err
}
