package usecases

import (
	"fmt"
	"os"
	"uploader/internal/constants"
	implements "uploader/pkg/aws"
)

func UploadFileChatUsecase(fileUUID string, extension string, privacy string, awsSdk *implements.AwsSdkImplementation) error {

	fileType := "image"

	if extension == "mp4" {
		fileType = "video"
	}

	bucket := os.Getenv("AWS_BUCKET_PUBLIC")

	if privacy == "private" {
		bucket = os.Getenv("AWS_BUCKET_PRIVATE")
	}

	//UPLOAD TO S3
	awsSdk.Upload(implements.UploadToS3{
		FileDir:     fmt.Sprintf("./tmp/%s/main.%s", fileUUID, extension),
		ContentType: fmt.Sprintf("%s/%s", fileType, extension),
		Bucket:      bucket,
		FileKey:     fileUUID,
		Folder:      constants.POST_FOLDER,
	})

	return nil
}
