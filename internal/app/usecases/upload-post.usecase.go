package usecases

import (
	"fmt"
	"os"
	"uploader/internal/constants"
	implements "uploader/pkg/aws"
	"uploader/utils"
)

func UploadPostUsecase(fileUUID string, extension string, privacy string, awsSdk *implements.AwsSdkImplementation) error {

	watemark := utils.NewWatemark(fileUUID, extension)
	fileType := "image"

	if extension == "mp4" {
		fileType = "video"
		err := watemark.WatemarkVideo()
		if err != nil {
			return err

		}

	} else {
		err := watemark.WatemarkImage()
		if err != nil {
			return err
		}
	}

	bucket := os.Getenv("AWS_BUCKET_PUBLIC")

	if privacy == "private" {
		bucket = os.Getenv("AWS_BUCKET_PRIVATE")
	}

	//UPLOAD TO S3
	err := awsSdk.Upload(implements.UploadToS3{
		FileDir:     fmt.Sprintf("./tmp/%s/watermark.%s", fileUUID, extension),
		ContentType: fmt.Sprintf("%s/%s", fileType, extension),
		Bucket:      bucket,
		FileKey:     fileUUID,
		Folder:      constants.POST_FOLDER,
	})

	if err != nil {
		return err
	}

	return nil
}
