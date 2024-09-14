package utils

import (
	"fmt"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func BlurImage(fileUUID string, extension string) error {

	// BLUR IMAGE

	err := ffmpeg_go.Input(fmt.Sprintf("./tmp/%s/compressed.%s", fileUUID, extension)).Filter("boxblur", ffmpeg_go.Args{"40", "40"}).
		Output(fmt.Sprintf("./tmp/%s/exposed.%s", fileUUID, extension)).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		fmt.Println(err, "Error blurring image")
		return err
	}

	return nil
}
