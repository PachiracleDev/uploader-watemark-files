package utils

import (
	"fmt"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func ShortVideo(fileUUID string, extension string, seconds int) error {

	err := ffmpeg_go.Input(fmt.Sprintf("./tmp/%s/compressed.%s", fileUUID, extension)).
		Output(fmt.Sprintf("./tmp/%s/exposed.%s", fileUUID, extension),
			ffmpeg_go.KwArgs{
				"t":   "5",    // Duración en segundos
				"c:v": "copy", // Copiar el stream de video sin recodificar
				"c:a": "copy", // Copiar el stream de audio sin recodificar
			}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()

	if err != nil {
		fmt.Println(err, "Error shortening video")
		return err
	}

	return nil

}
