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

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(conf.AWS.Region))

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
		fmt.Println("Error opening file", err)
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
		fmt.Println("Error uploading file to S3", err)
		return err
	}

	return nil
}
