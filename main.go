package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"uploader-image/config"
	"uploader-image/utils"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Sub        string `json:"sub"`
	Role       string `json:"role"`
	Preference string `json:"preference"`
	jwt.RegisteredClaims
}

func main() {
	e := echo.New()

	config.LoadEnv()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
	e.Use(echojwt.WithConfig(config))

	// Routes
	e.POST("/upload/:folder/:typeFile/:privacy", uploadFile)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

func getFile(size int64, fileType string) (string, error) {

	fileTypes := map[string]string{
		"image": "jpeg",
		"video": "mp4",
	}

	typeFile := fileTypes[fileType]

	//LIMITE DE 10MB
	if size > 10*1024*1024 {
		return "", fmt.Errorf("El archivo es demasiado grande")
	}

	if typeFile == "" {
		return "", fmt.Errorf("El tipo de archivo no es compatible")
	}

	return typeFile, nil
}

func uploadFile(c echo.Context) error {

	//
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	userId := claims.Sub

	// Obtener el tipo de archivo

	file, err := c.FormFile("file")

	if file == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "No se ha enviado ningún archivo",
		})
	}

	folder := c.Param("folder")

	if folder != "post" && folder != "users" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "El folder no es valido",
		})
	}

	size := file.Size

	typeFile, err := getFile(size, c.Param("typeFile"))
	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	privacy := c.Param("privacy")

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	//OBTENER UUID
	uuidGenerated := utils.GenerateUUID()

	if folder == "users" {
		uuidGenerated = userId
	}

	tempFile, err := os.CreateTemp("", fmt.Sprintf("file-%s.%s", uuidGenerated, typeFile))
	if err != nil {
		return err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, src)
	if err != nil {
		return err
	}

	input := tempFile.Name()

	output := fmt.Sprintf("file-watemark-%s.%s", uuidGenerated, typeFile)
	if folder == "post" {

		watemark := utils.NewWatemark(input, output)

		if typeFile == "jpeg" {
			err := watemark.WatemarkImage()
			if err != nil {
				return err
			}
		}

		if typeFile == "mp4" {
			err := watemark.WatemarkVideo()
			if err != nil {
				return err

			}
		}
	}

	if folder == "users" {
		output = uuidGenerated
		var wg sync.WaitGroup
		wg.Add(2)
		go utils.ResizeImage(input, output, "48", &wg)

		go utils.ResizeImage(input, output, "150", &wg)

		wg.Wait()

	}

	// Eliminar la imagen temporal del servidor
	err = os.Remove(tempFile.Name())
	if err != nil {
		return err
	}

	// Subir a S3
	uploaderS3 := utils.NewUploadToS3(output, uuidGenerated, privacy, typeFile)

	if folder == "post" {
		err = uploaderS3.UploadPost()
		if err != nil {
			return err
		}
	}

	if folder == "users" {
		err = uploaderS3.UploadProfile()
		if err != nil {
			return err
		}

	}

	if folder == "post" {
		err = os.Remove(output)
		if err != nil {
			return err
		}
	}

	if folder == "users" {
		err = os.Remove(fmt.Sprintf("%s_48px.jpeg", output))
		if err != nil {
			return err
		}

		err = os.Remove(fmt.Sprintf("%s_150px.jpeg", output))
		if err != nil {
			return err
		}
	}

	//json
	return c.JSON(http.StatusOK, echo.Map{
		"uuid": uuidGenerated,
	})

}
