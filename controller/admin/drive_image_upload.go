package admin

import (
	"context"
	"net/http"
	"portfolio/controller/auth"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func HandleAuth(c *gin.Context) {
	url := auth.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"AuthURL": url,
	})
}


func HandleCreateFolder(c *gin.Context) {
	parentFolderID := c.PostForm("parent_folder_id")
	folderName := c.PostForm("folder_name")

	client := auth.GetClient(auth.Config)
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrieve Drive client: %v", err)
		return
	}

	folder := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentFolderID},
	}
	folder, err = srv.Files.Create(folder).Do()
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to create folder: %v", err)
		return
	}
	c.String(http.StatusOK, "Folder created: %s\n", folder.Id)
}

func HandleUploadFile(c *gin.Context) {
	folderID := c.PostForm("folder_id")
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "File not found in request: %v", err)
		return
	}

	client := auth.GetClient(auth.Config)
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrieve Drive client: %v", err)
		return
	}

	f, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to open file: %v", err)
		return
	}
	defer f.Close()

	driveFile := &drive.File{
		Name:    file.Filename,
		Parents: []string{folderID},
	}

	driveFile, err = srv.Files.Create(driveFile).Media(f).Do()
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to upload file: %v", err)
		return
	}
	c.String(http.StatusOK, "File uploaded: %s\n", driveFile.Id)
}
