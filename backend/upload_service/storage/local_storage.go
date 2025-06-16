package storage

import (
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SaveToLocal(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	storagePath := "./storage/videos"
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		err := os.MkdirAll(storagePath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	filePath := fmt.Sprintf("%s/%d_%s", storagePath, time.Now().Unix(), file.Filename)

	// Use Fiber context to save file
	err := c.SaveFile(file, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
