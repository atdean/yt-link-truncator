package main

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func main() {
	// sourceUrl := "https://www.youtube.com/watch?v=YhxKbuO93ZI&t=3m&something=stuff"
	sourceUrl := "https://youtu.be/YhxKbuO93ZI?t=187"

	parsedURL, err := url.Parse(sourceUrl)
	if err != nil {
		os.Exit(1)
	}

	var videoCode string
	isFullLenHostName, err := regexp.MatchString(`[w*\.]*youtube\.com`, parsedURL.Host)
	if isFullLenHostName {
		queryComponents := parsedURL.Query()
		videoCode = queryComponents.Get("v")
		if videoCode == "" {
			os.Exit(2)
		}
	} else if parsedURL.Host == "youtu.be" {
		videoCode = parsedURL.EscapedPath()
		videoCode = strings.Trim(videoCode, "/")
	}

	reassembledURL := fmt.Sprintf("https://youtu.be/%s", videoCode)

	fmt.Println(reassembledURL)
}