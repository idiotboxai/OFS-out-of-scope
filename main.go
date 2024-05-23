package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

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
 / _ \ | |_  / __|  
| (_) || __| \__ \  
 \___/ |_|   |___/  
                    
`
	fmt.Print(color.RedString(logo))
}

// Function to check if a domain is in the out-of-scope list
func isOutOfScope(domain string, outOfScope []string) bool {
	for _, rule := range outOfScope {
		if strings.HasPrefix(rule, "*.") {
			if strings.HasSuffix(domain, rule[1:]) {
				return true
			}
		} else if strings.Contains(rule, "/*") {
			if strings.HasPrefix(domain, rule[:len(rule)-2]) {
				return true
			}
		} else {
			if domain == rule {
				return true
			}
		}
	}
	return false
}

// Function to read lines from a file
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

// Function to write lines to a file
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

func main() {
	
	subdomainsFile := flag.String("i", "", "Path to the subdomains file")
	outOfScopeFile := flag.String("s", "", "Path to the out-of-scope file")
	outputFile := flag.String("o", "", "Path to the output file")
	flag.Parse()


	if *subdomainsFile == "" || *outOfScopeFile == "" || *outputFile == "" {
		fmt.Println("Usage: ofc -i <subdomains_file> -s <out_of_scope_file> -o <output_file>")
		return
	}


	printLogo()

	subdomains, err := readLines(*subdomainsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading subdomains file: %v\n", err)
		return
	}
	outOfScope, err := readLines(*outOfScopeFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading out-of-scope file: %v\n", err)
		return
	}


	var filteredSubdomains []string
	for _, subdomain := range subdomains {
		if isOutOfScope(subdomain, outOfScope) {
			fmt.Println(color.RedString("Removed: %s", subdomain))
		} else {
			fmt.Println(color.GreenString("Retained: %s", subdomain))
			filteredSubdomains = append(filteredSubdomains, subdomain)
		}
	}

	err = writeLines(filteredSubdomains, *outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
	}

	fmt.Println(color.GreenString("\nProcessing complete!"))
}
