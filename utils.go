package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/likexian/whois"
)

// getTimezoneOffset returns the base and offset of the
// current system timezone like "GMT+2" or "CET-3"
func getTimezoneOffset() string {
	t := time.Now()
	zone, offset := t.Zone()
	offset = offset / 60 / 60
	offS := zone + "+" + strconv.Itoa(offset)
	if offset < 0 {
		offS = zone + "-" + strconv.Itoa(offset)
	} else if offset == 0 {
		offS = zone
	}
	return offS
}

// getHosterMail returns the abuse contact email address from
// the hoster of a given ip using public whois records.
// Returns empty string in case whois failed to return such address.
// If there are multiple abuse mailaddresses found, it takes the last one
// to prevent ripe address usage.
func getHosterMail(ip string) string {
	whoisResult, err := whois.Whois(ip)
	if err != nil {
		return ""
	}
	var re = regexp.MustCompile(`(?mi)abuse(?:.*:\s*|.*')([A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4})`)
	matches := re.FindAllStringSubmatch(string(whoisResult[:]), -1)
	if len(matches) == 0 {
		return ""
	}

	useMail := ""

	for _, mail := range matches {
		if len(mail) == 2 {
			useMail = mail[1]
			if *debugMode {
				fmt.Printf("- NOTE: Found email address %v in whois record.\n", useMail)
			}
		}
	}
	return useMail
}
