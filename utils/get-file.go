package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

type ResponseFile struct {
	FileUUID  string
	Extension string
	Dir       string
	Error     error
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func GetFile(c *fiber.Ctx,
	limitSize int64,
	extensionPermited []string,
) ResponseFile {
	//FORM FILE
	file, err := c.FormFile("file")

	if err != nil {
		return ResponseFile{
			Error: fmt.Errorf("file not found"),
		}
	}

	//VERIFY SIZE

	if file.Size > limitSize {
		return ResponseFile{
			Error: fmt.Errorf("file too large"),
		}
	}

	// Obtener la extensión del archivo
	extension := filepath.Ext(file.Filename)[1:]

	//VERIFY EXTENSION
	if !contains(extensionPermited, extension) {
		return ResponseFile{
			Error: fmt.Errorf("file extension not permited"),
		}
	}

	//GENERATE UUID KEY
	fileUUID := GenerateUUID()

	//CREAR UNA CARPETA CON EL UUID
	dir := fmt.Sprintf("./tmp/%s", fileUUID)
	os.Mkdir(dir, 0755)

	fileKey := fmt.Sprintf("%s/main.%s", fileUUID, extension)
	fileDir := fmt.Sprintf("./tmp/%s", fileKey)

	c.SaveFile(file, fileDir)

	return ResponseFile{
		FileUUID:  fileUUID,
		Extension: extension,
		Dir:       dir,
		Error:     nil,
	}
}
