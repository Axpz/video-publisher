package douyin

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/axpz/video-publisher/internal/config"
	"github.com/playwright-community/playwright-go"
)

const douyinURL = "https://www.douyin.com/"

type Client struct {
	Config config.Config
}

func NewClient(cfg config.Config) *Client {
	return &Client{Config: cfg}
}

func (c *Client) Auth(ctx context.Context) error {
	sessionFile := c.Config.SessionFile

	if err := playwright.Install(); err != nil {
		return err
	}

	pw, err := playwright.Run()
	if err != nil {
		return err
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		return err
	}
	defer browser.Close()

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		StorageStatePath: playwright.String(sessionFile),
	})
	if err != nil {
		return err
	}
	page, err := context.NewPage()
	if err != nil {
		return err
	}

	if _, err = page.Goto(douyinURL); err != nil {
		return err
	}

	fmt.Println("Complete the login in your browser, then press Enter to continue")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadBytes('\n')

	state, err := context.StorageState()
	if err != nil {
		return err
	}

	file, err := os.Create(sessionFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(state); err != nil {
		return err
	}

	fmt.Printf("Login session saved to %s\n", sessionFile)
	return nil
}
