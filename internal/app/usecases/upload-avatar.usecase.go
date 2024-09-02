package usecases

import (
	"fmt"
	"os"
	"sync"

	"uploader/internal/constants"
	"uploader/utils"

	implements "uploader/pkg/aws"
)

func UploadAvatar(fileId string, extension string, awsSdk *implements.AwsSdkImplementation) error {

	sizes := map[string]int{
		"sm": 48,
		"md": 150,
	}

	var wg sync.WaitGroup

	wg.Add(len(sizes))

	// Resize image
	for _, pixels := range sizes {
		go func(pixel int) {
			defer wg.Done()
			utils.ResizeImage(fileId, extension, pixel)
		}(pixels)
	}
	wg.Wait()

	var wg2 sync.WaitGroup

	wg2.Add(len(sizes))

	for size, pixels := range sizes {
		go func(pixel int, size string) {
			defer wg2.Done()
			awsSdk.Upload(implements.UploadToS3{
				FileDir:     fmt.Sprintf("./tmp/%s/resize_%d.%s", fileId, pixel, extension),
				ContentType: fmt.Sprintf("image/%s", extension),
				Bucket:      os.Getenv("AWS_BUCKET_PUBLIC"),
				FileKey:     fmt.Sprintf("%s_%s", fileId, size),
				Folder:      constants.USER_AVATAR_FOLDER,
			})
		}(pixels, size)
	}

	wg2.Wait()

	return nil
}
