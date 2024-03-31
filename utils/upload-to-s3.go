package utils

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type UploadToS3 struct {
	fileName string
	key      string
	privacy  string
	typeFile string
	svc      *s3.S3
}

func NewUploadToS3(fileName string, key string, privacy string, typeFile string) UploadToS3 {

	S3_UPLOAD_KEY := os.Getenv("S3_UPLOAD_KEY")
	S3_UPLOAD_SECRET := os.Getenv("S3_UPLOAD_SECRET")
	S3_UPLOAD_REGION := os.Getenv("S3_UPLOAD_REGION")

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(S3_UPLOAD_REGION),                                          // Cambia la región según tu necesidad
		Credentials: credentials.NewStaticCredentials(S3_UPLOAD_KEY, S3_UPLOAD_SECRET, ""), // Cambia "S3_UPLOAD_KEY" y "S3_UPLOAD_SECRET" con tus credenciales de AWS
	})

	// Crear un servicio S3
	svc := s3.New(sess)

	return UploadToS3{
		fileName: fileName,
		key:      key,
		privacy:  privacy,
		typeFile: typeFile,
		svc:      svc,
	}
}

func (u *UploadToS3) UploadProfile() error {
	var wgroup sync.WaitGroup
	wgroup.Add(2)

	go UploaderPhoto(fmt.Sprintf("%s_%spx.jpeg", u.fileName, "48"), fmt.Sprintf("%s-xs", u.key), u.svc, &wgroup)
	go UploaderPhoto(fmt.Sprintf("%s_%spx.jpeg", u.fileName, "150"), fmt.Sprintf("%s-sm", u.key), u.svc, &wgroup)

	wgroup.Wait()

	return nil
}

func UploaderPhoto(fileName string, key string, svc *s3.S3, wg *sync.WaitGroup) error {
	defer wg.Done()
	S3_UPLOAD_BUCKET := os.Getenv("S3_UPLOAD_BUCKET")
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Obtener información del archivo
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket:        aws.String(S3_UPLOAD_BUCKET),
		Key:           aws.String(fmt.Sprintf("/users/%s.jpeg", key)),
		Body:          file,
		ContentType:   aws.String("image/png"),
		ContentLength: aws.Int64(fileInfo.Size()),
	}

	// Subir el archivo a S3
	_, err = svc.PutObject(input)
	if err != nil {
		return err
	}

	return nil
}

func (u *UploadToS3) UploadPost() error {
	S3_UPLOAD_BUCKET := os.Getenv("S3_UPLOAD_BUCKET")

	if u.privacy == "private" {
		S3_UPLOAD_BUCKET = os.Getenv("S3_UPLOAD_BUCKET_PRIVATE")
	}

	// Abrir el archivo
	file, err := os.Open(u.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Obtener información del archivo
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	contentType := "image/jpeg"

	if u.typeFile == "mp4" {
		contentType = "video/mp4"
	}

	input := &s3.PutObjectInput{
		Bucket:        aws.String(S3_UPLOAD_BUCKET),
		Key:           aws.String(fmt.Sprintf("/post/%s.%s", u.key, u.typeFile)),
		Body:          file,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(fileInfo.Size()),
	}

	// Subir el archivo a S3
	_, err = u.svc.PutObject(input)
	if err != nil {
		return err
	}

	return nil
}
