package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type URLData struct {
	Host string
	Path string
	FullURL string
	ParameterSize int
	ParameterCount int
}

type Config struct {
	quiet bool
	filename string
	paramCount bool
	minParamCount int
}

func readFile(config Config) {

	file, err := os.Open(config.filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error while opening file! Make sure file exists")
		return
	}

	readAll := bufio.NewScanner(file)

	processedURLS := make(map[string]URLData)
	prevURL := URLData{}

	for readAll.Scan() {
		currentLine := readAll.Text()
		currentURL, parseError := url.Parse(currentLine)

		if parseError != nil {
			continue
		}

		parameters, _ := url.ParseQuery(currentURL.RawQuery)

		if processedURLS[currentURL.Host] == (URLData{}) {
			processedURLS[currentURL.Host] = URLData{
				Host:     currentURL.Host,
				Path:     currentURL.RawQuery,
				FullURL:  currentURL.Scheme+"://"+ currentURL.Host + currentURL.Path + currentURL.RawQuery,
				ParameterSize: len(currentURL.RawQuery),
				ParameterCount : len(parameters),
			}
		}

		if config.paramCount {
			if prevURL.Host == processedURLS[currentURL.Host].Host && prevURL.ParameterCount > processedURLS[currentURL.Host].ParameterCount {
				processedURLS[currentURL.Host] = prevURL
			}
		}else {
			if prevURL.Host == processedURLS[currentURL.Host].Host && prevURL.ParameterSize > processedURLS[currentURL.Host].ParameterSize {
				processedURLS[currentURL.Host] = prevURL
			}
		}

		prevURL = processedURLS[currentURL.Host]
	}

	for _, data := range processedURLS {
		if data.ParameterCount < config.minParamCount {
			continue
		}

		if config.quiet {
			fmt.Println(data.FullURL)
			continue
		}
		fmt.Println(fmt.Sprintf("Host: %s\nPath: %s\n", data.Host, data.FullURL))
	}
}

func main() {
	flag.Usage = func() {
		help := []string{
			"Filter given URLs based on query size",
			"",
			"Options:",
			"  -f,  --file			File to process",
			"  -c,  --count			Run analysis based on query count, otherwise query length",
			"  -mc, --min-count		Minimum parameter count",
			"  -q,  --quiet			Only print fullpath",
			"",
		}
		_, _ = fmt.Fprintf(os.Stderr, strings.Join(help, "\n"))
	}

	mConfig := Config{}

	flag.StringVar(&mConfig.filename, "file", "", "")
	flag.StringVar(&mConfig.filename, "f", "", "")

	flag.BoolVar(&mConfig.quiet, "quiet", false, "")
	flag.BoolVar(&mConfig.quiet, "q", false, "")

	flag.BoolVar(&mConfig.paramCount, "count", false, "")
	flag.BoolVar(&mConfig.paramCount, "c", false, "")

	flag.IntVar(&mConfig.minParamCount, "min-count", 1, "")
	flag.IntVar(&mConfig.minParamCount, "mc", 1, "")

	flag.Parse()

	if mConfig.filename == "" {
		fmt.Println("Please provide a file")
		return
	}

	readFile(mConfig)
}