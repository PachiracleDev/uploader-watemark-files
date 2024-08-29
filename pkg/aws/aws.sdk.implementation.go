package aws

import (
	"context"
	"fmt"
	"os"
	appConfig "uploader/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsSdkImplementation struct {
	s3Client *s3.Client
}

func NewSDKImplementation(conf *appConfig.Config) (*AwsSdkImplementation, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))

	if err != nil {
		return &AwsSdkImplementation{}, fmt.Errorf("error creating AWS session: %w", err)
	}

	return &AwsSdkImplementation{
		s3Client: s3.NewFromConfig(cfg),
	}, nil
}

type UploadToS3 struct {
	FileDir     string
	ContentType string
	Bucket      string
	FileKey     string
	Folder      string
}

func (u *AwsSdkImplementation) Upload(dto UploadToS3) error {

	file, err := os.Open(dto.FileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket:      aws.String(dto.Bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s", dto.Folder, dto.FileKey)),
		Body:        file,
		ContentType: aws.String(dto.ContentType),
	}

	_, err = u.s3Client.PutObject(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}
