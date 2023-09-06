package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

func searchWaybackMachine(target string, showDate bool) {
	apiURL := fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=txt&fl=original,timestamp", target)
	response, err := http.Get(apiURL)

	if err != nil {
		fmt.Printf("Error accessing the Wayback Machine for target: %s\n", target)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error reading the Wayback Machine response for target: %s\n", target)
			return
		}

		results := strings.Split(string(body), "\n")
		uniqueURLs := make(map[string]bool)
		var dateURLList []struct {
			Timestamp string
			URL       string
		}

		if len(results) > 1 {
			for _, result := range results[1:] {
				parts := strings.Fields(result)
				if len(parts) == 2 {
					url, timestamp := parts[0], parts[1]
					if !uniqueURLs[url] {
						uniqueURLs[url] = true
						dateURLList = append(dateURLList, struct {
							Timestamp string
							URL       string
						}{timestamp, url})
					}
				}
			}

			// Sort the list by modification date
			sort.Slice(dateURLList, func(i, j int) bool {
				return dateURLList[i].Timestamp < dateURLList[j].Timestamp
			})

			for _, entry := range dateURLList {
				if showDate {
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
		fmt.Printf("Error accessing the Wayback Machine for target: %s\n", target)
	}
}

func main() {
	targets := flag.String("targets", "", "Comma-separated list of targets to search (e.g., example.com,example2.com)")
	showDate := flag.Bool("details", false, "Show details from Web Archive")
	flag.Parse()

	if *targets == "" {
		fmt.Println(`
_  _  _ _______ __   __ ______  _______ _______ _     _    ______  _____ 
|  |  | |_____|   \_/   |_____] |_____| |       |____/    |  ____ |     |
|__|__| |     |    |    |_____] |     | |_____  |    \_ . |_____| |_____|
                                                           Version: 1.0.0
                                                    Author: @andersonbosa

USAGE: wayback.go [-details] -targets target1,target2,target3
`)
		os.Exit(1)
	}

	targetsList := strings.Split(*targets, ",")
	for _, target := range targetsList {
		searchWaybackMachine(target, *showDate)
	}
}
