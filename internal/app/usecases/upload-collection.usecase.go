package usecases

import (
	"fmt"
	"os"
	"uploader/internal/constants"
	"uploader/internal/presenters/dtos"
	implements "uploader/pkg/aws"
	"uploader/utils"
)

func UploadCollectionFileUsecase(dto dtos.UploadCollectionFile, fileUUID string, extension string, awsSdk *implements.AwsSdkImplementation) error {

	watemark := utils.NewWatemark(fileUUID, extension, fmt.Sprintf("gootitas.com/u/%s", dto.Username))
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

	//UPLOAD TO S3
	err := awsSdk.Upload(implements.UploadToS3{
		FileDir:     fmt.Sprintf("./tmp/%s/compressed.%s", fileUUID, extension),
		ContentType: fmt.Sprintf("%s/%s", fileType, extension),
		Bucket:      os.Getenv("AWS_BUCKET_PRIVATE"),
		FileKey:     fmt.Sprintf("%s/%s", dto.CollectionId, fileUUID),
		Folder:      constants.COLLECTION_FOLDER,
	})

	if err != nil {
		return err
	}

	//ELIMINAR CARPETA

	go os.RemoveAll(fmt.Sprintf("./tmp/%s", fileUUID))

	var body = map[string]string{
		"collectionId": dto.CollectionId,
		"tokenUpload":  dto.TokenUpload,
		"contentKey":   fileUUID,
		"contentType":  fileType,
	}

	// MANDAR REQUEST A API PARA ACTUALIZAR LA COLECCION
	utils.SendRequestHTTP(
		"http://localhost:9012/collections/insert-multimedia",
		"POST",
		map[string]string{
			"Content-Type": "application/json",
		},
		body,
	)

	return nil
}
