package usecases

import (
	"fmt"
	"os"

	"uploader/internal/constants"
	"uploader/utils"

	implements "uploader/pkg/aws"
)

func UploadBanner(fileId string, extension string, awsSdk *implements.AwsSdkImplementation) error {

	err := utils.CompressBanner(fileId, extension)

	if err != nil {
		return err
	}

	erro := awsSdk.Upload(implements.UploadToS3{
		FileDir:     fmt.Sprintf("./tmp/%s/compressed.%s", fileId, extension),
		ContentType: fmt.Sprintf("image/%s", extension),
		Bucket:      os.Getenv("AWS_BUCKET_PUBLIC"),
		FileKey:     fileId,
		Folder:      constants.USER_BANNER_FOLDER,
	})

	if erro != nil {
		return erro
	}

	return nil
}
