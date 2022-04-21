package scan

import (
	"bufio"
	"context"
	"coscanner/internal"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type DomainScanner struct{}

func (scanner *DomainScanner) Scan(flags internal.Flags, wg sync.WaitGroup) {
	if flags.WordList != "" {
		if file, err := os.Open(flags.WordList); err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				resolve(scanner.Text(), flags.Domain, flags.Verbose, wg)
			}
		} else {
			log.Fatal(err)
		}
	}

	if flags.Subdomain != "" {
		// single subdomain is filled
		for _, s := range strings.Split(flags.Subdomain, ",") {
			resolve(s, flags.Domain, flags.Verbose, wg)
		}
	}
}

func resolve(subdomain string, addr string, printerr bool, wg sync.WaitGroup) {
	wg.Add(1)
	go func() {
		domain := fmt.Sprintf("%s.%s", subdomain, addr)
		ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", domain)

		if err != nil {
			if printerr {
				printResult(domain, "NA")
			}
			wg.Done()
			return
		}

		printResult(domain, ips[0].String())
		wg.Done()
	}()
}

func printResult(domain string, ip string) {
	fmt.Printf("# [\033[32m%s\033[0m] \033[31m%s\033[0m\n", domain, ip)
}
