package html2png

import (
	"context"
	"github.com/playwright-community/playwright-go"
	"os"
	"strings"
)

var (
	DefaultChromePaths = []string{
		"/usr/bin/chromium-browser",
		"/usr/bin/chromium",
		"/usr/bin/google-chrome-stable",
		"/usr/bin/google-chrome",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"C:/Program Files (x86)/Google/Chrome/Application/chrome.exe",
		"C:/Program Files/Google/Chrome/Application/chrome.exe"}
)

func getChromePath() string {
	for _, path := range DefaultChromePaths {

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return path
		}
	}

	return ""
}

func HtmlToPng(ctx context.Context, html string, height int, width int) ([]byte, error) {
	// write val to html in temp
	// convert to png
	tempDir := os.TempDir()
	htmlFile := strings.Replace(tempDir+"\\temp.html", "\\", "/", -1)
	f, err := os.Create(htmlFile)
	if err != nil {
		return nil, err
	}
	
	_, err = f.WriteString(html)
	if err != nil {
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}

	client, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	// screenshot
	page, err := client.NewPage()
	if err != nil {
		return nil, err
	}

	_, err = page.Goto("file:///" + htmlFile)
	if err != nil {
		return nil, err
	}

	err = page.SetViewportSize(width, height)
	if err != nil {
		return nil, err
	}

	screenshot, err := page.Screenshot()
	if err != nil {
		return nil, err
	}

	_ = os.Remove(htmlFile)
	return screenshot, nil
}
