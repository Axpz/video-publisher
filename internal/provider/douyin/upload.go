package douyin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/playwright-community/playwright-go"
)

const uploadUrl = "https://creator.douyin.com/creator-micro/content/upload"
const uploadSelector = ".container-drag-upload-tL99XD"
const publishButtonSelector = "text=发布"

func (c *Client) Upload(ctx context.Context, filePath, metadataPath string) (string, error) {
	sessionFile := c.Config.SessionFile

	if _, err := os.Stat(sessionFile); err != nil {
		return "", fmt.Errorf("%s does not exist, please run `auth` first", sessionFile)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("Failed to get absolute path of video file: %v", err)
	}
	if _, err := os.Stat(absPath); err != nil {
		return "", fmt.Errorf("Video file does not exist: %s", absPath)
	}

	if err := playwright.Install(); err != nil {
		return "", err
	}

	pw, err := playwright.Run()
	if err != nil {
		return "", err
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		return "", err
	}
	defer browser.Close()

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		StorageStatePath: playwright.String(sessionFile),
	})
	if err != nil {
		return "", err
	}
	page, err := context.NewPage()
	if err != nil {
		return "", err
	}

	if _, err = page.Goto(uploadUrl); err != nil {
		return "", err
	}

	uploadTrigger := page.Locator(uploadSelector)
	if err := uploadTrigger.WaitFor(); err != nil {
		return "", err
	}

	time.Sleep(3 * time.Second)

	fileChooser, err := page.ExpectFileChooser(func() error {
		return uploadTrigger.Click()
	})
	if err != nil {
		return "", err
	}

	time.Sleep(3 * time.Second)

	if err := fileChooser.SetFiles([]string{absPath}); err != nil {
		return "", err
	}

	time.Sleep(3 * time.Second)

	publishButton := page.GetByText(publishButtonSelector)
	if err := publishButton.WaitFor(); err != nil {
		return "", err
	}

	if err := publishButton.Click(); err != nil {
		return "", err
	}

	time.Sleep(10 * time.Second)
	// Currently Douyin does not provide a simple permanent URL that can be obtained immediately,
	// so we return an empty string for now.
	return "", nil
}
