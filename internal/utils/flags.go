package utils

import (
	"strconv"
	"strings"
)

//Flags parameters
type Flags struct {
	Type      string //the type of scan
	Verbose   bool   //verbosity
	Domain    string //domain to scan
	Subdomain string //subdomain list to search
	WordList  string //the list of words to search subdomains
	Format    string //the output format
	Port      []int  //the ports to scan
}

func Parse(args []string) Flags {
	var result Flags

	parseArg(args, []string{"type", "t"}, true, &result.Type, extractString, "port")
	parseArg(args, []string{"verbose", "v"}, false, &result.Verbose, extractBool, false)
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
			if len(a) <= 1 || !strings.HasPrefix(a, flagPrefix) {
				continue
			}

			trim := TrimPrefix(a, flagPrefix)
			if cut, v, found := strings.Cut(trim, flagSeparator); cut == n {
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

func extractString(s string) string {
	return s
}

func extractBool(s string) bool {
	return true
}
