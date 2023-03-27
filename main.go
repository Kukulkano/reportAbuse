package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type detection struct {
	ip        string
	timestamp string
	words     []string
}

var detections map[string]detection
var debugMode *bool
var configFile *string

func main() {
	fmt.Println("Running reportAbuse")

	debugMode = flag.Bool("debug", false, "if available we are debugging")
	configFile = flag.String("config", "", "path to config json file")

	flag.Parse()

	if *debugMode {
		fmt.Println("DEBUG MODE. Will not send any emails to the hosters!")
	}

	loadConfig() // load and check configuration

	loadPatterns() // load detection patterns

	initDatabase() // init the database

	detections = make(map[string]detection)

	for _, logFile := range cfg.Logfiles {
		fmt.Println("Examine file " + logFile)
		examineFile(logFile)
	}

	closeDatabase()
}

func examineFile(logFile string) {
	// read the file
	content, err := os.ReadFile(logFile)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("- File size: %v bytes\n", len(content))

	timeZoneOffset := " " + getTimezoneOffset()

	// find all entries that match the "bad" words
	re := regexp.MustCompile(cfg.RegEx)
	matches := re.FindAllStringSubmatch(string(content[:]), -1)
	for _, match := range matches {
		tim := match[cfg.RegGroupDate]
		ip := match[cfg.RegGroupIP]
		page := match[cfg.RegGroupPage]
		if cfg.Mode == "direct" || cfg.Mode == "" {
			// take all entries found by regex
			if entry, ok := detections[ip]; ok {
				// add detected attack
				entry.words = append(entry.words, page)
				detections[ip] = entry
			} else {
				// create new entry
				var nEntry detection
				nEntry.timestamp = tim + timeZoneOffset
				nEntry.words = append(entry.words, page)
				nEntry.ip = ip
				detections[ip] = nEntry
			}
		}
		if cfg.Mode == "page" {
			// find all entries that match the "bad" words
			for _, word := range patterns {
				if strings.Contains(page, "/"+word) {
					if entry, ok := detections[ip]; ok {
						// add detected word
						entry.words = append(entry.words, page)
						detections[ip] = entry
					} else {
						// create new entry
						var nEntry detection
						nEntry.timestamp = tim + timeZoneOffset
						nEntry.words = append(nEntry.words, page)
						nEntry.ip = ip
						detections[ip] = nEntry
					}
					break
				}
			}
		}

	}

	// thin the processing list using minAttacks config value
	for key, detection := range detections {
		if len(detection.words) < cfg.MinAttacks {
			if *debugMode {
				fmt.Printf("- NOTE: Skip detected attack from IP %v because of not enough attacks (minAttacks).\n", key)
			}
			delete(detections, key)
		}
	}

	// Check against database to avoid multiple reports for the same attacker IP
	knownCount := 0
	for key, detection := range detections {
		val, _ := ipDB.Get([]byte(key))
		if len(val) == 0 {
			// ip not yet found in database
			fmt.Println("- Found new attacker from", key)
			hosterMail := getHosterMail(key)
			if !*debugMode {
				err = ipDB.Put([]byte(key), []byte(hosterMail))
				if err != nil {
					panic("Failed to add entry to database!")
				}
			} else {
				fmt.Printf("- NOTE: Do not remember IP %v because of debug mode.\n", key)
			}
			if hosterMail == "" {
				fmt.Println("- Cannot send email because I can't find abuse email")
				continue // with next entry
			}
			fmt.Println("- Sending email to", hosterMail, "   Attack timestamp", detection.timestamp)
			notifyHoster(detection, hosterMail)
		} else {
			// ip found, remove from processing list
			knownCount++
			delete(detections, key)
		}
	}

	fmt.Printf("Detected %v new attacks and %v already known attackers\n", len(detections), knownCount)
}
