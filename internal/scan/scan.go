package scan

import (
	"coscanner/internal"
	"sync"
)

const (
	ScanTypePort   string = "port"
	ScanTypeDomain string = "domain"
)

var scannerMap = map[string]Scanner{
	ScanTypePort:   &PortScanner{},
	ScanTypeDomain: &DomainScanner{},
}

type Scanner interface {
	Scan(flags internal.Flags, wg sync.WaitGroup)
}

// NewScanner Creates a new scan based on the command line flags
func NewScanner(flags internal.Flags) Scanner {
	if scanner, found := scannerMap[flags.Type]; found {
		return scanner
	}

	return nil
}
