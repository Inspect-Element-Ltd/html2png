package html2png

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
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

	pngFile := tempDir + "\\temp.png"

	_, err = f.WriteString(html)
	if err != nil {
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	args := []string{
		"--headless",
		"--no-sandbox",
		"--disable-crash-reporter",
		"--hide-scrollbars",
		"--default-background-color=00000000",
		"--disable-gpu",
		"--window-size=" + fmt.Sprintf("%d,%d", width, height),
		"--screenshot=" + pngFile,
		"file://" + htmlFile,
	}

	if err := exec.CommandContext(ctx, getChromePath(), args...).Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.New("takes screenshot got timeout")
		}

		return nil, err
	}

	if _, err := os.Stat(pngFile); os.IsNotExist(err) {
		return nil, err
	}

	png, err := os.ReadFile(pngFile)
	if err != nil {
		return nil, err
	}

	_ = os.Remove(htmlFile)
	_ = os.Remove(pngFile)

	return png, nil
}
