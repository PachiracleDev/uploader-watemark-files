package utils

import (
	"fmt"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type Watemark struct {
	input  string
	output string
}

func NewWatemark(input string, output string) Watemark {
	return Watemark{
		input:  input,
		output: output,
	}
}

func (w *Watemark) WatemarkImage() error {
	err := ffmpeg_go.Filter(
		[]*ffmpeg_go.Stream{
			ffmpeg_go.Input(w.input).Filter("scale", ffmpeg_go.Args{"600:-1"}),
			ffmpeg_go.Input("./images/logo.png").Filter("scale", ffmpeg_go.Args{"200:-1"}),
		}, "overlay", ffmpeg_go.Args{"10:10"}).
		Output(fmt.Sprintf("./%s", w.output)).OverWriteOutput().ErrorToStdOut().Run()

	return err
}

func (w *Watemark) WatemarkVideo() error {
	overlay := ffmpeg_go.Input("./images/logo.png").Filter("scale", ffmpeg_go.Args{"200:-1"})

	// Ejecutar ffmpeg con el filtro overlay modificado
	err := ffmpeg_go.Filter(
		[]*ffmpeg_go.Stream{
			ffmpeg_go.Input(w.input),
			overlay,
		}, "overlay", ffmpeg_go.Args{"10:10"}, ffmpeg_go.KwArgs{"enable": "gte(t,0)"}).
		Output(fmt.Sprintf("./%s", w.output)).OverWriteOutput().ErrorToStdOut().Run()

	return err
}
