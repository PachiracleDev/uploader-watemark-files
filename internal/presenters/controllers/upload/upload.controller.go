package controllers

import (
	"fmt"
	"os"
	"uploader/config"
	"uploader/internal/app/usecases"
	"uploader/internal/constants"
	"uploader/internal/presenters/dtos"
	"uploader/utils"

	"uploader/pkg/http"
	"uploader/pkg/validator"

	"github.com/gofiber/fiber/v2"

	implements "uploader/pkg/aws"
	//UTILS
)

func UploadController(
	http *http.HttpServer,
	conf *config.Config,

	validate *validator.XValidator,
	awsSdk *implements.AwsSdkImplementation,

) error {

	api := http.Group("/upload")

	// api.Use(http.AuthMiddleware())

	api.Post("/avatar", func(c *fiber.Ctx) error {

		result := utils.GetFile(
			c,
			constants.LIMIT_AVATAR_SIZE,
			[]string{"jpeg", "jpg", "png"},
		)

		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				map[string]interface{}{
					"message": result.Error.Error(),
				})
		}

		//USECASE
		usecases.UploadAvatar(result.FileUUID, result.Extension, awsSdk)

		//ELIMINAR CARPETA
		go os.RemoveAll(result.Dir)

		return c.JSON(map[string]interface{}{
			"fileKey": result.FileUUID,
		})

	})

	api.Post("/post", func(c *fiber.Ctx) error {

		fmt.Println("xd")
		// QUERY PARAMS PARSE
		uploadFileDto := new(dtos.UploadFile)
		if err := c.QueryParser(uploadFileDto); err != nil {
			return err
		}

		// Validation
		if errs := validate.Validate(uploadFileDto); errs != nil {
			return errs
		}

		result := utils.GetFile(
			c,
			constants.LIMIT_POST_SIZE,
			[]string{"jpeg", "jpg", "png", "mp4"},
		)

		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				map[string]interface{}{
					"message": result.Error.Error(),
				})
		}

		//USECASE
		usecases.UploadPostUsecase(result.FileUUID, result.Extension, uploadFileDto.Privacy, awsSdk)

		//ELIMINAR CARPETA
		go os.RemoveAll(result.Dir)

		return c.JSON(map[string]interface{}{
			"fileKey": result.FileUUID,
		})

	})

	api.Post("/chat", func(c *fiber.Ctx) error {

		// QUERY PARAMS PARSE
		uploadFileDto := new(dtos.UploadFile)
		if err := c.QueryParser(uploadFileDto); err != nil {
			return err
		}

		// Validation
		if errs := validate.Validate(uploadFileDto); errs != nil {
			return errs
		}

		result := utils.GetFile(
			c,
			constants.LIMIT_POST_SIZE,
			[]string{"jpeg", "jpg", "png", "mp4"},
		)

		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				map[string]interface{}{
					"message": result.Error.Error(),
				})
		}

		//USECASE
		usecases.UploadFileChatUsecase(result.FileUUID, result.Extension, uploadFileDto.Privacy, awsSdk)

		go os.RemoveAll(result.Dir)

		return c.JSON(map[string]interface{}{
			"fileKey": result.FileUUID,
		})

	})

	return nil
}
