package middleware

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func UploadFile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		file, err := c.FormFile("image")
		if err != nil && c.Request().Method == "PATCH" {
			c.Set("dataFile", "false")
			return next(c)
		}

		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, "Error Retrieving the File")
		}

		const MAX_UPLOAD_SIZE = 10 << 20
		if c.Request().ContentLength > MAX_UPLOAD_SIZE {
			return c.JSON(http.StatusBadRequest, Result{Code: http.StatusBadRequest, Message: "Max size in 10mb"})
		}

		tempFile, err := os.CreateTemp("uploads", "waysbeans-*.png")
		if err != nil {
			fmt.Println(err)
			fmt.Println("path upload error")
			return c.JSON(http.StatusInternalServerError, err)
		}
		defer tempFile.Close()

		src, err := file.Open()
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		defer src.Close()

		if _, err = io.Copy(tempFile, src); err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		data := tempFile.Name()

		c.Set("dataFile", data)
		return next(c)
	}
}
