package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// Create Echo instance
	e := echo.New()

	// Routes
	e.PUT("/api/v1/files", putFile)
	e.POST("/api/v1/cid", getCid)
	e.POST("/api/v1/move", moveFile)
	e.DELETE("/api/v1/delete", deleteFile)
	e.POST("/api/v1/files", list)

	// Starting Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Fatal(e.Start(":" + port))
}

/*
	 Request: PUT /api/v1/putFile
			{
				"fileName": "test.txt",
				"path": "/test.txt", // Default is
				"cid": "QmW2WQi7j6c7UgJTarActp7tDNikE4B2qXtFCfLPdsgaTQ" // CID of the file
				"size": "123456" // Size of the file in bytes
			}
		Return
			{
				"id": "5f9e1b9b9c9c2b0007e1b5f1",
				"fileName": "test.txt",
				"path": "/test.txt", // Default is
			}
*/
func putFile(c echo.Context) error {
	db := Connect()
	// convert size to int64
	size, e := strconv.ParseInt(c.FormValue("size"), 10, 64)
	if e != nil {
		log.Println(e)
		return c.JSON(http.StatusBadRequest, e.Error())
	}

	newFile := File{
		Path:     c.FormValue("path"),
		FileName: c.FormValue("fileName"),
		CID:      c.FormValue("cid"),
		Size:     size,
	}

	err := db.Create(&newFile).Error
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, newFile)
}

/*
	 Request: POST /api/v1/move
			{
				"sourcePath": "/test.txt",
				"destinationPath": "/test2.txt"
			}
		Return
			{
				"id": "5f9e1b9b9c9c2b0007e1b5f1",
				"fileName": "test.txt",
				"path": "/test.txt", // Default is
			}
*/
func moveFile(c echo.Context) error {
	db := Connect()

	sourcePath := c.FormValue("sourcePath")
	destinationPath := c.FormValue("destinationPath")
	
	var file File
	if err := db.Where("path = ?", sourcePath).First(&file).Error; err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	file.Path = destinationPath
	if err := db.Save(&file).Error; err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, "File moved successfully")
}

/*
	 Request: DELETE /api/v1/file
			{
				"path": "/test.txt"
			}
		Return
			{
				"id": "5f9e1b9b9c9c2b0007e1b5f1",
				"fileName": "test.txt",
				"path": "/test.txt", // Default is
			}
*/
func deleteFile(c echo.Context) error {
	db := Connect()

	path := c.FormValue("path")

	if err := db.Where("path = ?", path).Delete(&File{}).Error; err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, "File deleted successfully")
}

/*
	 Request: POST /api/v1/list
			{
				"path": "/test.txt"
			}
		Return
			{
				[
					{
				"id": "5f9e1b9b9c9c2b0007e1b5f1",
				"fileName": "test.txt",
				"path": "/test.txt", // Default is
				{
				"id": "5f9e1b9b9c9c2b0007e1asd1",
				"fileName": "tes2t.txt",
				"path": "/test2.txt", // Default is
				]
			}
*/
func list(c echo.Context) error {
	db := Connect()

	path := c.FormValue("path")

	// Find all files with path matching the regex and return their paths in an array
	var files []File

	if err := db.Raw("SELECT * FROM files WHERE path LIKE ? || '%' AND deleted_at IS NULL", path).Scan(&files).Error; err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var paths []string
	for _, file := range files {
		paths = append(paths, file.Path)
	}

	return c.JSON(http.StatusOK, files)
}

func getCid(c echo.Context) error {
	db := Connect()

	path := c.FormValue("path")

	var file File
	if err := db.Where("path = ?", path).First(&file).Error; err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, file.CID)
}
