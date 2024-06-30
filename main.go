package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
)

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	Reset = "\033[0m"
)

func printLogo() {
	logo := `
    )   (     (
 ( /(   )\ )  )\ )
 )\()) (()/( (()/(
((_)\   /(_)) /(_))
  ((_) (_))_|(_))
 / * \ | |*  / __|
| (_) || **| \** \
 \___/ |_|   |___/
`
	fmt.Print(color.RedString(logo))
}

func matchesDomain(domain string, rule string) bool {
	rule = strings.TrimSpace(rule)
	if strings.HasPrefix(rule, "*.") {
		wildcard := rule[2:]
		return strings.HasSuffix(domain, wildcard)
	} else if strings.HasSuffix(rule, "/*") {
		baseDomain := rule[:len(rule)-2]
		return domain == baseDomain || strings.HasPrefix(domain, baseDomain+"/")
	} else {
		return domain == rule
	}
}

func isAllowed(domain string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, rule := range allowed {
		if matchesDomain(domain, rule) {
			return true
		}
	}
	return false
}

func isDisallowed(domain string, disallowed []string) bool {
	for _, rule := range disallowed {
		fmt.Printf("Checking domain '%s' against rule '%s'\n", domain, rule)
		if matchesDomain(domain, rule) {
			fmt.Printf("Match found: '%s' is disallowed by rule '%s'\n", domain, rule)
			return true
		}
	}
	return false
}

func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(lines []string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func filterDomains(subdomains []string, allowed []string, disallowed []string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, subdomain := range subdomains {
		if isDisallowed(subdomain, disallowed) {
			fmt.Println(color.RedString("Removed (Disallowed): %s", subdomain))
		} else if !isAllowed(subdomain, allowed) {
			fmt.Println(color.RedString("Removed (Not Allowed): %s", subdomain))
		} else {
			fmt.Println(color.GreenString("Retained: %s", subdomain))
			results <- subdomain
		}
	}
}

func main() {
	subdomainsFile := flag.String("IL", "", "Path to the subdomains file")
	allowedDomains := flag.String("a", "", "Allowed domains (comma-separated)")
	disallowedDomains := flag.String("d", "", "Disallowed domains (comma-separated)")
	outputFile := flag.String("o", "", "Path to the output file")
	flag.Parse()

	if *subdomainsFile == "" || *outputFile == "" {
		fmt.Println("Usage: ofc -IL <subdomains_file> -a <allowed_domains> -d <disallowed_domains> -o <output_file>")
		return
	}

	printLogo()

	subdomains, err := readLines(*subdomainsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading subdomains file: %v\n", err)
		return
	}

	var allowed []string
	if *allowedDomains != "" {
		allowed = strings.Split(*allowedDomains, ",")
	}

	var disallowed []string
	if *disallowedDomains != "" {
		disallowed, err = readLines(*disallowedDomains)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading disallowed domains file: %v\n", err)
			return
		}
	}

	fmt.Println("Disallowed rules:")
	for _, rule := range disallowed {
		fmt.Printf("- %s\n", rule)
	}

	results := make(chan string, len(subdomains))
	var wg sync.WaitGroup

	numWorkers := 4
	chunkSize := (len(subdomains) + numWorkers - 1) / numWorkers

	for i := 0; i < len(subdomains); i += chunkSize {
		end := i + chunkSize
		if end > len(subdomains) {
			end = len(subdomains)
		}
		wg.Add(1)
		go filterDomains(subdomains[i:end], allowed, disallowed, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var filteredSubdomains []string
	for subdomain := range results {
		filteredSubdomains = append(filteredSubdomains, subdomain)
	}

	err = writeLines(filteredSubdomains, *outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		return
	}

	fmt.Println(color.GreenString("\nProcessing complete!"))
}
