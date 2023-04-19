package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
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
	e.POST("/api/v1/files/move", moveFile)
	e.DELETE("/api/v1/files", deleteFile)
	e.POST("/api/v1/files", list)

	// Starting Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Fatal(e.Start(":" + port))
}


/* Request: PUT /api/v1/putFile
		{
			"fileName": "test.txt",
			"path": "/test.txt", // Default is 
			"cid": "QmW2WQi7j6c7UgJTarActp7tDNikE4B2qXtFCfLPdsgaTQ" // CID of the file
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

	newFile := File{
		Path:     c.FormValue("path"),
		FileName: c.FormValue("fileName"),
		CID:      c.FormValue("CID"),
	}

	err := db.Create(&newFile).Error
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, newFile)
}

/* Request: POST /api/v1/move
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

/* Request: DELETE /api/v1/file
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

/* Request: POST /api/v1/list
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

	//A regular expression to match subpaths with one more level
	regex := "^" + regexp.QuoteMeta(path) + "/[^/]+$"

	var files []File
	if err := db.Where("path ~ ?", regex).Order("path asc").Find(&files).Error; err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, files)
}
