package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

var (
	Config *oauth2.Config
)

func HandleCallback(c *gin.Context) {
	code := c.Query("code")
	tok, err := Config.Exchange(context.Background(), code)
	if err != nil {
		c.String(http.StatusBadRequest, "Unable to retrieve token from web: %v", err)
		return
	}
	SaveToken("token.json", tok)
	c.String(http.StatusOK, "Authentication successful! You can now use the Drive API.")
}

func GoogleClientSetUp() {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatalf("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables must be set")
	}

	Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{drive.DriveFileScope, drive.DriveReadonlyScope},
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/callback",
	}

}
func GetClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := TokenFromFile(tokFile)
	if err != nil {
		log.Fatalf("Unable to read token from file: %v", err)
	}
	return config.Client(context.Background(), tok)
}

func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
