package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
)

var (
	desc    = "Google Transparencyreport subdomains finder"
	version = "0.9.0"
	author  = "LordEver (@Lord3ver)"

	out     = true
	outFile = ""
	rmWC    = true
	rmDup   = true
	rmExt   = false
	lookup  = true
)

// Parse the content and returns domain list and next page token
func cntParser(content []byte) (domains []string, nextPageToken string) {
	var dat [][]interface{}

	content = bytes.ReplaceAll(content, []byte("\n"), []byte(""))
	content = bytes.ReplaceAll(content, []byte(")]}'"), []byte(""))

	if err := json.Unmarshal(content, &dat); err != nil {
		log.Panic(err)
	}

	d := reflect.ValueOf(dat[0][1])
	for i := 0; i < d.Len(); i++ {
		slice, ok := d.Index(i).Interface().([]interface{})
		if !ok {
			panic("Value error")
		}
		domains = append(domains, slice[1].(string))
	}

	// Return domain list if this is the last page (next page token is nil)
	npt := reflect.ValueOf(dat[0][3]).Index(1).Interface()
	if npt == nil {
		return domains, ""
	}
	nextPageToken = reflect.ValueOf(dat[0][3]).Index(1).Interface().(string)

	return
}

// Get page body
func getPage(domain string, page string) []byte {
	var url string
	if page != "" {
		url = fmt.Sprintf("https://transparencyreport.google.com/transparencyreport/api/v3/httpsreport/ct/certsearch/page?p=%s", page)
	} else {
		url = fmt.Sprintf("https://transparencyreport.google.com/transparencyreport/api/v3/httpsreport/ct/certsearch?include_subdomains=true&domain=%s", domain)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 503 {
		return []byte("")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	return body
}

func dnsExist(domain string) bool {
	_, err := net.LookupIP(domain)
	if err != nil {
		return false
	}
	return true
}

// Get domains
func expose(domain string) []string {
	var domains []string
	var token string
	var domainsTmp []string

	for {
		domainsTmp, token = cntParser(getPage(domain, token))

		for _, d := range domainsTmp {
			// Avoid duplicates
			if rmDup && sliceContainsString(domains, d) {
				continue
			}
			// Avoid wildcard domains
			if rmWC {
				if strings.Index(d, "*") == 0 {
					continue
				}
			}
			// Avoid external domains
			if rmExt {
				if strings.Index(d, domain) == -1 {
					continue
				}
			}
			if lookup {
				if !dnsExist(d) {
					//fmt.Println("DNS NOT FOUND:", d)
					continue
				}
			}

			if out {
				fmt.Println(d)
			}
			domains = append(domains, d)
		}
		domainsTmp = []string{}

		// Reached last page
		if token == "" {
			return domains
		}

	}

}

// sliceContainsString check if dlist contains domain
func sliceContainsString(dlist []string, domain string) bool {
	for _, d := range dlist {
		if d == domain {
			return true
		}
	}
	return false
}

// writeToFile save found domains to a file
func writeToFile(fn string, domains []string) {
	file, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range domains {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	file.Close()
}

func main() {

	banner := `
     ______     __  _____       __        __                      _           
    / ____/____/ /_/ ___/__  __/ /_  ____/ /___  ____ ___  ____ _(_)___  _____
   / / __/ ___/ __/\__ \/ / / / __ \/ __  / __ \/ __ '__ \/ __ '/ / __ \/ ___/
  / /_/ / /__/ /_ ___/ / /_/ / /_/ / /_/ / /_/ / / / / / / /_/ / / / / (__  ) 
  \____/\___/\__//____/\__,_/_.___/\__,_/\____/_/ /_/ /_/\__,_/_/_/ /_/____/  
																			  
	`

	fmt.Println(banner)
	fmt.Printf("\t%s\t\n\n\tVersion:\t%s\n\tAuthor:\t\t%s\n\n", desc, version, author)

	domainPtr := flag.String("d", "", "Target domain. E.g. bing.com")
	outPtr := flag.Bool("out", true, "Print results to stdout")
	outFilePtr := flag.String("outfile", "", "Specify an output file when completed. Create or append if exists.")
	rmWCPtr := flag.Bool("rmw", true, "Remove wildcard domains, ex. *.uber.com")
	rmDuplicatesPtr := flag.Bool("rmd", true, "Remove duplicates")
	rmExternalPtr := flag.Bool("rme", false, "Remove external domains, like xyz.com for uber.com domain (default false)")
	lookupPtr := flag.Bool("lookup", true, "Do DNS lookups for the domains to see which ones exist")

	flag.Parse()

	if *domainPtr == "" {
		fmt.Println("Domain is a mandatory parameter!")
		flag.PrintDefaults()
		return
	}
	if *outPtr == false && *outFilePtr == "" {
		fmt.Println("\"out\" is false and \"outfile\" is not set, where do you wanna to go?\n\nStdout output enabled.")
		*outPtr = true
	}
	out = *outPtr
	outFile = *outFilePtr
	rmWC = *rmWCPtr
	rmDup = *rmDuplicatesPtr
	rmExt = *rmExternalPtr
	lookup = *lookupPtr

	domains := expose(*domainPtr)

	if *outFilePtr != "" && len(domains) > 0 {
		writeToFile(*outFilePtr, domains)
	}

	fmt.Println("----------------------------------------")
	fmt.Printf("Total domains/subdomains: %v\n", len(domains))
	fmt.Println("----------------------------------------")

	fmt.Println("\nBye!")

}
