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
	Format    string
	Port      []int
}

func (c Flags) String() string {
	return fmt.Sprintf("[%v %v %v %v %v %v %v]", c.Type, c.Verbose, c.Domain, c.Subdomain, c.WordList, c.Format, c.Port)
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
	parseArg(args, []string{"format", "f"}, true, &result.Format, extractString, "text")

	// assume the last element of the args list is always the domain
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
	var portSlice []int

	for _, s := range strings.Split(ports, ",") {
		if lower, upper, found := strings.Cut(s, "-"); found {
			l, err := strconv.Atoi(lower)
			if err != nil {
				continue
			}
			u, err := strconv.Atoi(upper)
			if err != nil {
				continue
			}
			for i := l; i <= u; i++ {
				portSlice = append(portSlice, i)
			}
		} else {
			v, err := strconv.Atoi(s)
			if err != nil {
				continue
			}
			portSlice = append(portSlice, v)
		}
	}

	return portSlice
}
