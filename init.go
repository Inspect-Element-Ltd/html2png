package html2png

import "github.com/playwright-community/playwright-go"

func Init() {
	err := playwright.Install()
	if err != nil {
		panic(err)
	}
}
