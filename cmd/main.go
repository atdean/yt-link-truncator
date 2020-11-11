package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func main() {
	// Goal: allow input from either stdin, a CLI flag, or via HTTP interface
	// CLI flag+arg takes precedence over stdin
	// HTTP interface is started with a flag and will ignore stdin and CLI arg

	// sourceUrl := "https://www.youtube.com/watch?v=YhxKbuO93ZI&t=3m&something=stuff"
	// sourceUrl := "https://youtu.be/YhxKbuO93ZI?t=187"

	//convertedURL, err := parseURL(sourceUrl)
	//if err != nil {
	//	os.Exit(1)
	//}
	//fmt.Println(convertedURL)
}

func fetchSourceURL() (string, error) {
	var sourceUrl string

	if !func(sourceUrl *string) bool {
		// Parse command line flags (if any)
		flag.StringVar(sourceUrl, "url", "", "The Youtube URL to shorten.")
		flag.Parse()

		// If URL CLI flag is unset, check for piped string input
		if *sourceUrl == "" {
			fmt.Println("No URL flag provided... checking stdin")
			fi, err := os.Stdin.Stat()
			// AND char device flag w/ file mode to check if we got a string piped in via stdin
			if err == nil && (fi.Mode() & os.ModeCharDevice) == 0 {
				if inputBytes, err := ioutil.ReadAll(os.Stdin); err == nil {
					*sourceUrl = strings.TrimSpace(string(inputBytes))
				}
			}
		}

		return len(*sourceUrl) > 0
	}(&sourceUrl) {
		return "", errors.New("no URL was provided")
	} else {
		fmt.Println(sourceUrl)
		return sourceUrl, nil
	}
}

func parseURL(sourceUrl string) (string, error) {
	parsedURL, err := url.Parse(sourceUrl)
	if err != nil {
		return "", nil
	}

	var videoCode string
	isFullLenHostName, err := regexp.MatchString(`[w*\.]*youtube\.com`, parsedURL.Host)
	if isFullLenHostName {
		queryComponents := parsedURL.Query()
		videoCode = queryComponents.Get("v")
		if videoCode == "" {
			return "", errors.New("invalid Youtube URL")
		}
	} else if parsedURL.Host == "youtu.be" {
		videoCode = parsedURL.EscapedPath()
		videoCode = strings.Trim(videoCode, "/")
	}

	reassembledURL := fmt.Sprintf("https://youtu.be/%s", videoCode)

	return reassembledURL, nil
}