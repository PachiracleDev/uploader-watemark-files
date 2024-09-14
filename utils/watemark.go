package utils

import (
	"fmt"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type Watemark struct {
	fileUUID  string
	extension string
	text      string
}

func NewWatemark(fileUUID string, extension string, text string) Watemark {
	return Watemark{
		fileUUID:  fileUUID,
		extension: extension,
		text:      text,
	}
}

func (w *Watemark) WatemarkImage() error {

	// CENTER WATERMARK
	// err := ffmpeg_go.Filter(
	// 	[]*ffmpeg_go.Stream{
	// 		ffmpeg_go.Input(fmt.Sprintf("./tmp/%s/main.%s", w.fileUUID, w.extension)).Filter("scale", ffmpeg_go.Args{"700:-1"}),
	// 		ffmpeg_go.Input("./images/logo.png").Filter("scale", ffmpeg_go.Args{"300:-1"}),
	// 	}, "overlay", ffmpeg_go.Args{"(main_w-overlay_w)/2:(main_h-overlay_h)/2"}).
	// 	Output(fmt.Sprintf("./tmp/%s/watermark.%s", w.fileUUID, w.extension)).OverWriteOutput().ErrorToStdOut().Run()

	// BOTTOM RIGHT WATERMARK
	err := ffmpeg_go.Filter(
		[]*ffmpeg_go.Stream{
			ffmpeg_go.Input(fmt.Sprintf("./tmp/%s/main.%s", w.fileUUID, w.extension)).Filter("scale", ffmpeg_go.Args{"700:-1"}),
			ffmpeg_go.Input("./images/logo.png").Filter("scale", ffmpeg_go.Args{"300:-1"}),
		}, "overlay", ffmpeg_go.Args{"main_w-overlay_w-10:main_h-overlay_h-10"}).
		Output(fmt.Sprintf("./tmp/%s/watermark.%s", w.fileUUID, w.extension)).OverWriteOutput().ErrorToStdOut().Run()

	return err
}

func (w *Watemark) WatemarkText() error {

	// BOTTOM RIGHT WATERMARK TEXT
	err := ffmpeg_go.Filter(
		[]*ffmpeg_go.Stream{
			ffmpeg_go.Input(fmt.Sprintf("./tmp/%s/main.%s", w.fileUUID, w.extension)),
		}, "drawtext", ffmpeg_go.Args{fmt.Sprintf("fontfile=Arial.ttf:text='%s':fontcolor=white:fontsize=24:x=w-tw-10:y=h-th-10", w.text)}).
		Output(fmt.Sprintf("./tmp/%s/watermark.%s", w.fileUUID, w.extension)).OverWriteOutput().ErrorToStdOut().Run()

	return err
}

func (w *Watemark) WatermarkVideoText() error {

	err := ffmpeg_go.Filter(
		[]*ffmpeg_go.Stream{
			ffmpeg_go.Input(fmt.Sprintf("./tmp/%s/main.%s", w.fileUUID, w.extension)),
		}, "drawtext", ffmpeg_go.Args{fmt.Sprintf("fontfile=Arial.ttf:text='%s':fontcolor=white:fontsize=24:x=w-tw-10:y=h-th-10", w.text)}).
		Output(
			fmt.Sprintf("./tmp/%s/compressed.%s", w.fileUUID, w.extension),
			ffmpeg_go.KwArgs{
				"c:v":    "libx264", // Especifica el codec de video
				"c:a":    "aac",     // Especifica el codec de audio
				"map":    "0:a",     // Mapea el audio original
				"vcodec": "libx264", // Esto hace que el video se comprima
				"crf":    "28",      // Control de tasa constante para compresión

			},
		).OverWriteOutput().ErrorToStdOut().Run()

	return err

}

func (w *Watemark) WatemarkVideo() error {
	overlay := ffmpeg_go.Input("./images/logo.png").Filter("scale", ffmpeg_go.Args{"200:-1"})

	err := ffmpeg_go.Filter(
		[]*ffmpeg_go.Stream{
			ffmpeg_go.Input(fmt.Sprintf("./tmp/%s/main.%s", w.fileUUID, w.extension)),
			overlay,
		}, "overlay", ffmpeg_go.Args{"main_w-overlay_w-10:main_h-overlay_h-10"}, ffmpeg_go.KwArgs{"enable": "gte(t,0)"}).
		Output(
			fmt.Sprintf("./tmp/%s/watermark.%s", w.fileUUID, w.extension),
			ffmpeg_go.KwArgs{
				"c:v": "libx264", // Especifica el codec de video
				"c:a": "aac",     // Especifica el codec de audio (puedes usar "copy" para no recodificar)
				"map": "0:a",     // Mapea el audio original
			},
		).OverWriteOutput().ErrorToStdOut().Run()

	return err
}
