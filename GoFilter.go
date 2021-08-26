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
	Host           string
	Path           string
	FullURL        string
	ParameterSize  int
	ParameterCount int
}

type Config struct {
	quiet         bool
	filename      string
	forLength     bool
	minParamCount int
	readFromSTDIN bool
}

func processFile(config Config) {

	var file *os.File
	var readAll *bufio.Scanner
	var err error

	if !config.readFromSTDIN {
		file, err = os.Open(config.filename)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error while opening file! Make sure file exists")
			return
		}
		readAll = bufio.NewScanner(file)
	} else {
		readAll = bufio.NewScanner(os.Stdin)
	}

	processedURLS := make(map[string]URLData)

	for readAll.Scan() {
		currentLine := readAll.Text()
		currentURL, parseError := url.Parse(currentLine)

		if parseError != nil {
			continue
		}
		currentPath := currentURL.Host + currentURL.Path
		parameters, _ := url.ParseQuery(currentURL.RawQuery)

		currentURLData := URLData{
			Host:           currentURL.Host,
			Path:           currentURL.RawQuery,
			FullURL:        currentLine,
			ParameterSize:  len(currentURL.RawQuery),
			ParameterCount: len(parameters),
		}

		if processedURLS[currentPath] == (URLData{}) {
			processedURLS[currentPath] = currentURLData
			continue
		}

		if config.forLength && processedURLS[currentPath].ParameterSize < len(currentURL.RawQuery) {
			processedURLS[currentPath] = currentURLData
		} else if !config.forLength && processedURLS[currentPath].ParameterCount < len(parameters) {
			processedURLS[currentPath] = currentURLData
		}
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
			"  -f,  --file          File to process",
			"  -s   --std           Read from standard input",
			"  -l,  --length        Run analysis based on query length",
			"  -mc, --min-count     Minimum parameter count",
			"  -q,  --quiet         Only print fullpath",
			"",
		}
		_, _ = fmt.Fprintf(os.Stderr, strings.Join(help, "\n"))
	}

	mConfig := Config{}

	flag.StringVar(&mConfig.filename, "file", "", "")
	flag.StringVar(&mConfig.filename, "f", "", "")

	flag.BoolVar(&mConfig.readFromSTDIN, "s", false, "")
	flag.BoolVar(&mConfig.readFromSTDIN, "std", false, "")

	flag.BoolVar(&mConfig.quiet, "quiet", false, "")
	flag.BoolVar(&mConfig.quiet, "q", false, "")

	flag.BoolVar(&mConfig.forLength, "length", false, "")
	flag.BoolVar(&mConfig.forLength, "l", false, "")

	flag.IntVar(&mConfig.minParamCount, "min-count", 1, "")
	flag.IntVar(&mConfig.minParamCount, "mc", 1, "")

	flag.Parse()

	if mConfig.filename == "" && !mConfig.readFromSTDIN {
		fmt.Println("Please provide a file")
		return
	}

	processFile(mConfig)
}
