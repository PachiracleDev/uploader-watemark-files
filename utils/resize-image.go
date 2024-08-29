package utils

import (
	"fmt"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func ResizeImage(fileId string, extension string, pixels string) error {

	err := ffmpeg_go.Input(fmt.Sprintf("./tmp/%s/main.%s", fileId, extension)).Filter("scale", ffmpeg_go.Args{fmt.Sprintf("%s:-1", pixels)}).
		Output(fmt.Sprintf("./tmp/%s/%s.%s", fileId, pixels, extension)).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return err
	}

	return nil
}
