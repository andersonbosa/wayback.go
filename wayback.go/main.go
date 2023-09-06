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

func searchWaybackMachine(domain string, showDate bool) {
	apiURL := fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=txt&fl=original,timestamp", domain)
	response, err := http.Get(apiURL)

	if err != nil {
		fmt.Printf("Error accessing the Wayback Machine for domain: %s\n", domain)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error reading the Wayback Machine response for domain: %s\n", domain)
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
					fmt.Printf("Modification Date for domain %s: %s\n", domain, entry.Timestamp)
					fmt.Printf("URL: %s\n", entry.URL)
					fmt.Printf("Link to Web Archive: http://web.archive.org/web/%s/%s\n\n", entry.Timestamp, entry.URL)
				} else {
					fmt.Println(entry.URL)
				}
			}
		} else {
			fmt.Printf("No captures found for domain: %s\n", domain)
		}
	} else {
		fmt.Printf("Error accessing the Wayback Machine for domain: %s\n", domain)
	}
}

func main() {
	domains := flag.String("domains", "", "Comma-separated list of domains to search (e.g., example.com,example2.com)")
	showDate := flag.Bool("details", false, "Show modification date and link to Web Archive")
	flag.Parse()

	if *domains == "" {
		fmt.Println("Usage: wayback.go [-details] -domains domain1,domain2,domain3")
		os.Exit(1)
	}

	domainList := strings.Split(*domains, ",")
	for _, domain := range domainList {
		searchWaybackMachine(domain, *showDate)
	}
}
