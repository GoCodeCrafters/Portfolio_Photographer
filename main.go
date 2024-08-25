package main

import (
	"encoding/base64"
	"fmt"
	adminpkg "portfolio/controller/admin"
	"portfolio/controller/auth"
	clientpkg "portfolio/controller/client"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	//google auth setup
	auth.GoogleClientSetUp()

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"base64Encode": func(b []byte) string {
			return base64.StdEncoding.EncodeToString(b)
		},
	})

	// Grouping routes under "/admin"
	admin := r.Group("/admin")
	{
		admin.GET("/", adminpkg.HandleAuth)
		admin.GET("/callback", auth.HandleCallback)
		admin.POST("/create_folder", adminpkg.HandleCreateFolder)
		admin.POST("/upload_file", adminpkg.HandleUploadFile)
	}
	// Grouping routes under "/fetch"
	fetch := r.Group("/fetch")
	{
		fetch.GET("/:folderID", clientpkg.HandleFetchImages)
	}

	r.LoadHTMLGlob("templates/*")
	r.Run(":8080")
}
