package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"portfolio/controller/auth"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Image struct {
	Data     []byte
	MimeType string
}
// image fetching from drive with a specific floder id {} 
func HandleFetchImages(c *gin.Context) {
	folderID := c.Param("folderID")
	client := auth.GetClient(auth.Config)

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrieve Drive client: %v", err)
		return
	}

	query := fmt.Sprintf("'%s' in parents and mimeType contains 'image/'", folderID)
	fileList, err := srv.Files.List().Q(query).Fields("files(id, mimeType)").Do()
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrieve files: %v", err)
		return
	}

	var images []Image
	for _, file := range fileList.Files {
		log.Println("Processing file", file.Name)

		// Read the image file content
		fileContent, err := srv.Files.Get(file.Id).Download()
		if err != nil {
			log.Println("Error downloading file:", err)
			continue
		}
		defer fileContent.Body.Close()

		imageData, err := io.ReadAll(fileContent.Body)
		if err != nil {
			log.Println("Error reading file content:", err)
			continue
		}

		log.Println("Image data length:", len(imageData))

		image := Image{
			Data:     imageData,
			MimeType: file.MimeType,
		}
		images = append(images, image)
	}

	// Pass the images to the HTML file
	c.HTML(http.StatusOK, "client.html", gin.H{
		"images": images,
	})
}
