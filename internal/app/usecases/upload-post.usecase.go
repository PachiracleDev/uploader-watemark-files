package usecases

import (
	"fmt"
	"os"
	"uploader/internal/constants"
	implements "uploader/pkg/aws"
	"uploader/utils"
)

func UploadPostUsecase(fileUUID string, extension string, privacy string, username string, awsSdk *implements.AwsSdkImplementation) error {

	watemark := utils.NewWatemark(fileUUID, extension, fmt.Sprintf("gootitas.com/u/%s", username))
	fileType := "image"

	if extension == "mp4" {
		fileType = "video"
		err := watemark.WatermarkVideoText()
		if err != nil {
			return err

		}
	} else {
		err := watemark.WatemarkText()
		if err != nil {
			return err
		}
	}

	//COMPRESS IMAGE
	if extension != "mp4" {
		err := utils.CompressImage(fileUUID, extension)
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
		FileDir:     fmt.Sprintf("./tmp/%s/compressed.%s", fileUUID, extension),
		ContentType: fmt.Sprintf("%s/%s", fileType, extension),
		Bucket:      bucket,
		FileKey:     fileUUID,
		Folder:      constants.POST_FOLDER,
	})

	if err != nil {
		return err
	}

	if privacy == "private" {
		if extension == "mp4" {
			err := utils.ShortVideo(fileUUID, extension, 5)
			if err != nil {
				return err
			}
		} else {
			err := utils.BlurImage(fileUUID, extension)
			if err != nil {
				return err
			}
		}

		err := awsSdk.Upload(implements.UploadToS3{
			FileDir:     fmt.Sprintf("./tmp/%s/exposed.%s", fileUUID, extension),
			ContentType: fmt.Sprintf("%s/%s", fileType, extension),
			Bucket:      os.Getenv("AWS_BUCKET_PUBLIC"),
			FileKey:     fmt.Sprintf("%s_exposed", fileUUID),
			Folder:      constants.POST_FOLDER,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
