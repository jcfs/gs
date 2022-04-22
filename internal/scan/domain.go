package scan

import (
	"bufio"
	"context"
	"fmt"
	"gs/internal/utils"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type DomainScanner struct{}

type SubdomainStatus struct {
	Subdomain string `json:"subdomain"`
	Ip        string `json:"ip,omitempty"`
	Status    string `json:"status"`
}

type DomainScanResult struct {
	Elapsed    float64           `json:"elapsed"`
	Subdomains []SubdomainStatus `json:"subdomains,omitempty"`
}

func (scanner *DomainScanner) Scan(flags utils.Flags, wg *sync.WaitGroup) Result {
	var (
		chwg       sync.WaitGroup
		statusChan = make(chan SubdomainStatus)
		resultChan = make(chan DomainScanResult)
		done       = make(chan int)
	)

	go out(statusChan, resultChan, done, wg)

	if flags.WordList != "" {
		if file, err := os.Open(flags.WordList); err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				resolve(scanner.Text(), flags.Domain, &chwg, statusChan)
			}
		} else {
			log.Fatal(err)
		}
	}

	if flags.Subdomain != "" {
		// single subdomain is filled
		for _, s := range strings.Split(flags.Subdomain, ",") {
			resolve(s, flags.Domain, &chwg, statusChan)
		}
	}

	chwg.Wait()
	done <- 1
	return <-resultChan
}

func resolve(subdomain string, addr string, chwg *sync.WaitGroup, statusChan chan SubdomainStatus) {
	chwg.Add(1)
	go func() {
		domain := fmt.Sprintf("%s.%s", subdomain, addr)
		if ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", domain); err != nil {
			statusChan <- SubdomainStatus{Subdomain: domain, Status: utils.DomainNotFound}
		} else {
			statusChan <- SubdomainStatus{Subdomain: domain, Ip: ips[0].String(), Status: utils.DomainFound}
		}
		chwg.Done()
	}()
}

func out(statusChan chan SubdomainStatus, resultChan chan DomainScanResult, done chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	var (
		subdomains []SubdomainStatus
		start      = time.Now()
	)
	for {
		select {
		case r := <-statusChan:
			subdomains = append(subdomains, r)
		case <-done:
			resultChan <- DomainScanResult{Subdomains: subdomains, Elapsed: time.Since(start).Seconds()}
			wg.Done()
		}
	}
}
