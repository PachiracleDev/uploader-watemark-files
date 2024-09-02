package usecases

import (
	"fmt"
	"os"
	"uploader/internal/constants"
	implements "uploader/pkg/aws"
)

func UploadVouchersUsecase(fileUUID string, extension string, awsSdk *implements.AwsSdkImplementation) error {

	fileType := "image"

	bucket := os.Getenv("AWS_BUCKET_PUBLIC")

	// UPLOAD TO S3
	err := awsSdk.Upload(implements.UploadToS3{
		FileDir:     fmt.Sprintf("./tmp/%s/main.%s", fileUUID, extension),
		ContentType: fmt.Sprintf("%s/%s", fileType, extension),
		Bucket:      bucket,
		FileKey:     fileUUID,
		Folder:      constants.VOUCHER_FOLDER,
	})

	if err != nil {
		return err
	}

	return nil
}
