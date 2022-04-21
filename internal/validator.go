package internal

import (
	"coscanner/internal/scan"
	"errors"
	"fmt"
)

const usage = `
Usage: gs [options...] <domain>
 -t, --type <type>      The scan type (domain, port, ...)
 -v, --verbose          Also print unresolvable domains
 -o, --output <file>    Write output to file
DOMAIN:
 -s, --subdomain <data> The subdomain to test (ie: www,ns1,cloud)
 -w, --wordlist <file>  The word list to use
PORT:
 -p, --port <ports>     The port(s) to scan (ex: "1", "1-10", "1,2,3")

`

func Validate(flags Flags) error {
	if err := validateFlags(flags); err != nil {
		fmt.Print(usage)
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

//validateScan validates the scan configuration for each scan type
func validateFlags(flags Flags) error {
	if flags.Domain == "" {
		return errors.New("domain must be present")
	}

	switch flags.Type {
	case scan.ScanTypePort:
		if len(flags.Port) == 0 {
			return errors.New("missing -p/--port option")
		}
	case scan.ScanTypeDomain:
		if flags.Subdomain == "" && flags.WordList == "" {
			return errors.New("missing -w/--wordlist or -s/--subdomain option")
		}
	default:
		return fmt.Errorf("%s is an invalid scan type", flags.Type)
	}

	return nil
}
