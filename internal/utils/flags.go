package utils

import (
	"fmt"
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

func (c Flags) String() string {
	return fmt.Sprintf("[%v %v %v %v %v %v]", c.Type, c.Verbose, c.Domain, c.Subdomain, c.WordList, c.Port)
}

func Parse(args []string) Flags {
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
			if cut, v, found := strings.Cut(trim, separator); cut == n {
				var value string

				if found {
					value = v
				} else if keyValue {
					value = args[i+1]
				}

				*ref = extract(value)
			}
		}
	}
}

//
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
