package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type cliFlags struct {
	sourceURL string
}

func main() {
	// Goal: allow input from either stdin, a CLI flag, or via HTTP interface
	// CLI flag+arg takes precedence over stdin
	// TODO :: ^^ Consider swapping that priority? stdin first, CLI second?
	// HTTP interface is started with a flag and will ignore stdin and CLI arg

	// sourceUrl := "https://www.youtube.com/watch?v=YhxKbuO93ZI&t=3m&something=stuff"
	// sourceUrl := "https://youtu.be/YhxKbuO93ZI?t=187"

	flags, err := parseCLIFlags()
	if err != nil {
		log.Fatal(err)
	}
	convertedURL, err := parseURL(flags.sourceURL)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(convertedURL)
}

func parseCLIFlags() (*cliFlags, error){
	flags := cliFlags{}

	sourceURL, err := fetchSourceURL()
	if err != nil || sourceURL == "" {
		return nil, errors.New("FATAL: Unable to parse a source URL.")
	} else {
		flags.sourceURL = sourceURL
	}

	return &flags, nil
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
		// TODO :: this needs to check for a watch?v= parameter
		videoCode = parsedURL.EscapedPath()
		videoCode = strings.Trim(videoCode, "/")
	}

	reassembledURL := fmt.Sprintf("https://youtu.be/%s", videoCode)

	return reassembledURL, nil
}