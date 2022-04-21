package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	prefix    = "-"
	separator = "="
)

//Flags name definition
type Flags struct {
	Type      string
	Verbose   bool
	Domain    string
	Subdomain string
	WordList  string
	Port      []int
}

func Parse() Flags {
	args := os.Args[1:]
	var result Flags

	extractString := func(s string) string {
		return s
	}

	isPresent := func(s string) bool {
		return true
	}

	parseArg(args, []string{"type", "t"}, true, &result.Type, extractString, "port")
	parseArg(args, []string{"verbose", "v"}, false, &result.Verbose, isPresent, false)
	parseArg(args, []string{"subdomain", "s"}, true, &result.Subdomain, extractString, "")
	parseArg(args, []string{"wordlist", "w"}, true, &result.WordList, extractString, "")
	parseArg(args, []string{"port", "p"}, true, &result.Port, parseScanPortFlags, GetCommonPorts())

	// assume the last element of the args list if always the domain

	if len(args) > 0 {
		result.Domain = args[len(args)-1]
	}
	return result
}

func parseArg[T any](args []string, names []string, keyValue bool, ref *T, extract func(string) T, def T) {
	*ref = def
	for _, n := range names {
		for i, a := range args {
			// we only care about option keys here
			if len(a) <= 1 || !strings.HasPrefix(a, prefix) {
				continue
			}

			trim := TrimPrefix(a, prefix)
			if n == trim {
				var value, found = "", false
				if _, value, found = strings.Cut(trim, separator); !found {
					if keyValue {
						value = args[i+1]
					}
				}
				*ref = extract(value)
			}
		}
	}
}

func parseScanPortFlags(ports string) []int {
	// check if it is a single port
	if port, err := strconv.Atoi(ports); err == nil {
		return []int{port}
	}

	// check if comma separated
	if strings.Contains(ports, ",") {
		var portSlice []int
		for _, i := range strings.Split(ports, ",") {
			port, err := strconv.Atoi(i)
			if err != nil {
				continue
			}
			portSlice = append(portSlice, port)
		}
		return portSlice
	}

	// parse %d-%d format ie, 1-100 and expand it to 1,2,3...100
	upper, lower := 0, 0
	if _, err := fmt.Sscanf(ports, "%d-%d", &lower, &upper); err == nil {
		var portSlice []int
		for i := lower; i <= upper; i++ {
			portSlice = append(portSlice, i)
		}
		return portSlice
	}

	return nil

}
