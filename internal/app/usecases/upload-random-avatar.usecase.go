package usecases

import (
	"fmt"
	"os"
	"sync"
	"time"

	//"uploader/internal/constants"
	"uploader/internal/constants"
	implements "uploader/pkg/aws"
	"uploader/utils"
)

func UploadRandomAvatar(awsSdk *implements.AwsSdkImplementation) (string, error) {

	extension := "png"

	// Generar URL aleatoria de DiceBear
	avatarUrl := fmt.Sprintf("https://api.dicebear.com/8.x/adventurer/png?seed=%s&width=500&height=500", fmt.Sprintf("avatar_%d", time.Now().UnixNano()))

	//GENERATE UUID KEY
	fileUUID := utils.GenerateUUID()

	//CREAR UNA CARPETA CON EL UUID
	dir := fmt.Sprintf("./tmp/%s", fileUUID)
	os.Mkdir(dir, 0755)

	fileKey := fmt.Sprintf("%s/main.%s", fileUUID, extension)
	fileDir := fmt.Sprintf("./tmp/%s", fileKey)

	// Descargar imagen
	if err := utils.DownloadImageFromURL(avatarUrl, fileDir); err != nil {
		return fileUUID, fmt.Errorf("error descargando avatar: %w", err)
	}

	// Resize image
	sizes := map[string]int{
		"sm": 54,
		"md": 250,
	}

	var wg sync.WaitGroup

	wg.Add(len(sizes))

	for _, pixels := range sizes {
		go func(pixel int) {
			defer wg.Done()
			utils.ResizeImage(fileUUID, extension, pixel)
		}(pixels)
	}
	wg.Wait()

	var wg2 sync.WaitGroup

	wg2.Add(len(sizes))

	for size, pixels := range sizes {
		go func(pixel int, size string) {
			defer wg2.Done()
			awsSdk.Upload(implements.UploadToS3{
				FileDir:     fmt.Sprintf("./tmp/%s/resize_%d.%s", fileUUID, pixel, extension),
				ContentType: fmt.Sprintf("image/%s", extension),
				Bucket:      os.Getenv("AWS_BUCKET_PUBLIC"),
				FileKey:     fmt.Sprintf("%s_%s", fileUUID, size),
				Folder:      constants.USER_AVATAR_FOLDER,
			})
		}(pixels, size)
	}

	wg2.Wait()

	return fileUUID, nil
}
