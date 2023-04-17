package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
			"size": 1234, // Size of the file in bytes
			"cid": "QmW2WQi7j6c7UgJTarActp7tDNikE4B2qXtFCfLPdsgaTQ" // CID of the file
		}
	Return
		{
			"id": "5f9e1b9b9c9c2b0007e1b5f1",
			"fileName": "test.txt",
			"path": "/test.txt", // Default is
			"size": 1234, // Size of the file in bytes
		}
*/
func putFile(c echo.Context) error {
	// Get Connect
	client, db := connect()
	defer client.Disconnect(c.Request().Context())

	collection := db.Collection(os.Getenv("MONGODB_COLLECTION_NAME"))

	fileSize,e := strconv.ParseInt(c.FormValue("size"), 10, 64)
	if e != nil {
		log.Println(e)
		return c.JSON(http.StatusBadRequest, e.Error())
	}


	newFile := FileSystemNode{
		Path: c.FormValue("path"),
		FileName: c.FormValue("fileName"),
		Size: fileSize,
		CID: c.FormValue("CID"),
	}

	

	_, err := collection.InsertOne(c.Request().Context(), newFile)
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
			"size": 1234, // Size of the file in bytes
		}
*/
func moveFile(c echo.Context) error {
	// Get Connect
	client, db := connect()
	defer client.Disconnect(c.Request().Context())

	collection := db.Collection(os.Getenv("MONGODB_COLLECTION_NAME"))

	sourcePath := c.FormValue("sourcePath")
	destinationPath := c.FormValue("destinationPath")

	filter := bson.D{{"path", sourcePath}}
	update := bson.D{
		{"$set", bson.D{
			{"path", destinationPath},
		}},
	}

	_, err := collection.UpdateOne(c.Request().Context(), filter, update)
	if err != nil {
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
			"size": 1234, // Size of the file in bytes
		}
*/
func deleteFile(c echo.Context) error {
	// Get Connect
	client, db := connect()
	defer client.Disconnect(c.Request().Context())

	collection := db.Collection(os.Getenv("MONGODB_COLLECTION_NAME"))

	path := c.FormValue("path")

	filter := bson.D{{"path", path}}

	_, err := collection.DeleteOne(c.Request().Context(), filter)
	if err != nil {
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
			"size": 1234, // Size of the file in bytes},
			{
			"id": "5f9e1b9b9c9c2b0007e1asd1",
			"fileName": "tes2t.txt",
			"path": "/test2.txt", // Default is
			"size": 1234, // Size of the file in bytes}
			]
		}
*/
func list(c echo.Context) error {
	// Connect to the database
	client, db := connect()
	defer client.Disconnect(c.Request().Context())

	collection := db.Collection(os.Getenv("MONGODB_COLLECTION_NAME"))

	path := c.FormValue("path")

	// Create a regular expression to match subpaths with one more level
	regex := "^" + regexp.QuoteMeta(path) + "/[^/]+$"
	filter := bson.D{{"path", bson.M{"$regex": regex}}}

	findOptions := options.Find().SetSort(bson.D{{"path", 1}}) // Sort the results by path

	cursor, err := collection.Find(c.Request().Context(), filter, findOptions)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	var files []FileSystemNode
	if err = cursor.All(c.Request().Context(), &files); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, files)
}

