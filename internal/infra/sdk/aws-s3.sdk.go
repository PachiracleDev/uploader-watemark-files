package sdk

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AwsS3Sdk struct {
	serviceS3 *s3.S3
}

func NewAwsS3Sdk(serviceS3 *s3.S3) AwsS3Sdk {
	S3_UPLOAD_KEY := os.Getenv("S3_UPLOAD_KEY")
	S3_UPLOAD_SECRET := os.Getenv("S3_UPLOAD_SECRET")
	S3_UPLOAD_REGION := os.Getenv("S3_UPLOAD_REGION")

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(S3_UPLOAD_REGION),                                          // Cambia la región según tu necesidad
		Credentials: credentials.NewStaticCredentials(S3_UPLOAD_KEY, S3_UPLOAD_SECRET, ""), // Cambia "S3_UPLOAD_KEY" y "S3_UPLOAD_SECRET" con tus credenciales de AWS
	})

	// Crear un servicio S3
	svc := s3.New(sess)

	return AwsS3Sdk{
		serviceS3: svc,
	}
}

type UploadToS3 struct {
	FileDir     string
	ContentType string
	Bucket      string
	typeFile    string
	FileKey     string
	Folder      string
}

func (u *AwsS3Sdk) Upload(dto UploadToS3) error {

	file, err := os.Open(dto.FileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket:        aws.String(dto.Bucket),
		Key:           aws.String(fmt.Sprintf("/%s/%s.%s", dto.Folder, dto.FileKey, dto.typeFile)),
		Body:          file,
		ContentType:   aws.String(dto.ContentType),
		ContentLength: aws.Int64(fileInfo.Size()),
	}

	_, err = u.serviceS3.PutObject(input)
	if err != nil {
		return err
	}

	return nil
}
