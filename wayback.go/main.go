package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
)

// Result represents a Wayback Machine capture result.
type Result struct {
	Timestamp string
	URL       string
}

func searchWaybackMachine(target string, showDetails bool) {
	apiURL := fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=txt&fl=original,timestamp", target)
	response, err := http.Get(apiURL)

	if err != nil {
		fmt.Printf("Error accessing the Wayback Machine for target: %s - %v\n", target, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error reading the Wayback Machine response for target: %s - %v\n", target, err)
			return
		}

		results := strings.Split(string(body), "\n")
		uniqueURLs := make(map[string]bool)
		var dateURLList []Result

		if len(results) > 1 {
			for _, result := range results[1:] {
				parts := strings.Fields(result)
				if len(parts) == 2 {
					url, timestamp := parts[0], parts[1]
					if !uniqueURLs[url] {
						uniqueURLs[url] = true
						dateURLList = append(dateURLList, Result{timestamp, url})
					}
				}
			}

			// Sort the list by modification date
			sort.Slice(dateURLList, func(i, j int) bool {
				return dateURLList[i].Timestamp < dateURLList[j].Timestamp
			})

			for _, entry := range dateURLList {
				if showDetails {
					fmt.Println("---")
					fmt.Printf("Target: %s\n", entry.URL)
					fmt.Printf("Last updated %s: %s\n", target, entry.Timestamp)
					fmt.Printf("Web Archive link: http://web.archive.org/web/%s/%s\n", entry.Timestamp, entry.URL)
				} else {
					fmt.Println(entry.URL)
				}
			}
		} else {
			fmt.Printf("No captures found for target: %s\n", target)
		}
	} else {
		fmt.Printf("Error accessing the Wayback Machine for target: %s - Status Code: %d\n", target, response.StatusCode)
	}
}

func printUsage() {
	fmt.Println(`
_  _  _ _______ __   __ ______  _______ _______ _     _    ______  _____ 
|  |  | |_____|   \_/   |_____] |_____| |       |____/    |  ____ |     |
|__|__| |     |    |    |_____] |     | |_____  |    \_ . |_____| |_____|
                                                           Version: 1.0.0
                                                    Author: @andersonbosa

USAGE: wayback.go [-details] -targets target1,target2,target3
`)
}

func main() {
	targets := flag.String("targets", "", "Comma-separated list of targets to search (e.g., example.com,example2.com)")
	showDetails := flag.Bool("details", false, "Show details from Web Archive")
	flag.Parse()

	if *targets == "" {
		printUsage()
		os.Exit(1)
	}

	targetsList := strings.Split(*targets, ",")
	for _, target := range targetsList {
		searchWaybackMachine(target, *showDetails)
	}
}
