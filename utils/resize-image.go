package utils

import (
	"fmt"
	"sync"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func ResizeImage(inputPath string, output string, pixels string, wg *sync.WaitGroup) error {
	// Redimensionar a 48px
	defer wg.Done()
	err := ffmpeg_go.Input(inputPath).Filter("scale", ffmpeg_go.Args{fmt.Sprintf("%s:-1", pixels)}).
		Output(fmt.Sprintf("%s_%spx.jpeg", output, pixels)).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return err
	}

	return nil
}
