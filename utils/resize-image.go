package utils

import (
	"fmt"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func ResizeImage(fileId string, extension string, pixels int) error {

	err := ffmpeg_go.Input(fmt.Sprintf("./tmp/%s/main.%s", fileId, extension)).Filter("scale", ffmpeg_go.Args{fmt.Sprintf("%d:-1", pixels)}).
		Output(fmt.Sprintf("./tmp/%s/resize_%d.%s", fileId, pixels, extension)).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		fmt.Println(err, "Error resizing image")
		return err
	}

	return nil
}

func CrobImage(fileId string, extension string, width int, height int) error {

	err := ffmpeg_go.Input(fmt.Sprintf("./tmp/%s/main.%s", fileId, extension)).
		Filter("crop", ffmpeg_go.Args{fmt.Sprintf("%d:%d", width, height)}). // Recortar la imagen
		Output(fmt.Sprintf("./tmp/%s/crop_%d_%d.%s", fileId, width, height, extension)).
		OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		fmt.Println(err, "Error cropping image")
		return err
	}

	return nil
}
