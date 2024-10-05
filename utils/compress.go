package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

func CompressVideo(fileKey string, format string) {

}

func CompressImage(fileKey string, format string) error {

	// Abrir archivo
	file, err := os.Open("./tmp/" + fileKey + "/watermark." + format)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Crear archivo de salida
	out, err := os.Create("./tmp/" + fileKey + "/compressed." + format)
	if err != nil {
		return err
	}
	defer out.Close()

	// Guardar la imagen según su formato
	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 80}) // Comprimir JPEG
	case "png":
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		return encoder.Encode(out, img) // Comprimir PNG
	default:
		return fmt.Errorf("formato de imagen no soportado: %s", format)
	}
}

func CompressBanner(fileKey string, format string) error {

	// Abrir archivo
	file, err := os.Open("./tmp/" + fileKey + "/main." + format)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Crear archivo de salida
	out, err := os.Create("./tmp/" + fileKey + "/compressed." + format)
	if err != nil {
		return err
	}
	defer out.Close()

	// Guardar la imagen según su formato
	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 80}) // Comprimir JPEG
	case "png":
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		return encoder.Encode(out, img) // Comprimir PNG
	default:
		return fmt.Errorf("formato de imagen no soportado: %s", format)
	}
}
