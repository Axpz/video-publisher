package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/axpz/video-publisher/internal/app"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	yt "google.golang.org/api/youtube/v3"
)

type Client struct {
	Config app.Config
}

func NewClient(cfg app.Config) *Client {
	return &Client{Config: cfg}
}

func (c *Client) Auth() error {
	_, err := c.httpClient(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("YouTube authorization completed successfully")
	return nil
}

func (c *Client) httpClient(ctx context.Context) (*http.Client, error) {
	secrets, err := os.ReadFile(c.Config.ClientSecretsFile)
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(secrets, yt.YoutubeUploadScope)
	if err != nil {
		return nil, err
	}

	tokenFile := c.Config.TokenFile
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		if err := saveToken(tokenFile, tok); err != nil {
			return nil, err
		}
	}
	return config.Client(ctx, tok), nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Please open this link in your browser to authorize the application: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Save token to file: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
