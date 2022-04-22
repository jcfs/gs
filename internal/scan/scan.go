package scan

import (
	"gs/internal/utils"
	"sync"
)

const (
	TypePort   string = "port"
	TypeDomain string = "domain"
)

var scannerMap = map[string]Scanner{
	TypePort:   &PortScanner{},
	TypeDomain: &DomainScanner{},
}

type Result interface {
}

type Scanner interface {
	Scan(flags utils.Flags, wg *sync.WaitGroup) Result
}

// NewScanner Creates a new scan based on the command line flags
func NewScanner(flags utils.Flags) Scanner {
	if scanner, found := scannerMap[flags.Type]; found {
		return scanner
	}
	return nil
}
