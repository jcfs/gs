package report

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gs/internal/scan"
	"gs/internal/utils"
)

type Reporter interface {
	Report(result scan.Result)
}

const (
	portScanHeaderFormat   = "%7s%7s%20s\n"
	portScanLineFormat     = "%7d%s%7s\033[0m%20s\n"
	domainScanHeaderFormat = "%26s%11s%17s\n"
	domainScanLineFormat   = "%26s%11s%17s\n"
)

type PortScanTextReporter struct {
	flags utils.Flags
}

type DomainScanTextReporter struct {
	flags utils.Flags
}

type JsonReporter struct {
	flags utils.Flags
}

type XmlReporter struct {
	flags utils.Flags
}

func (r PortScanTextReporter) Report(result scan.Result) {
	scanResult := result.(scan.PortScanResult)

	fmt.Printf(portScanHeaderFormat, "Port", "Status", "Service")
	fmt.Printf(portScanHeaderFormat, "------", "------", "-------------------")
	for _, v := range scanResult.Ports {
		color := getColorByStatus(v.Status)

		if r.flags.Verbose && v.Status != utils.PortOpen {
			fmt.Printf(portScanLineFormat, v.Port, color, v.Status, utils.GetPortDescription(v.Port))
		}

		if v.Status == utils.PortOpen {
			fmt.Printf(portScanLineFormat, v.Port, color, v.Status, utils.GetPortDescription(v.Port))
		}

	}
	fmt.Printf("\nScanned %d ports in %.2fs.\n", len(scanResult.Ports), scanResult.Elapsed)
}

func (r DomainScanTextReporter) Report(result scan.Result) {
	scanResult := result.(scan.DomainScanResult)

	fmt.Printf(domainScanHeaderFormat, "Domain", "Status", "Service")
	fmt.Printf(domainScanHeaderFormat, "-------------------------", "----------", "----------------")
	for _, v := range scanResult.Subdomains {
		if r.flags.Verbose && v.Status != utils.DomainFound {
			fmt.Printf(domainScanLineFormat, v.Subdomain, v.Status, v.Ip)
		}

		if v.Status == utils.DomainFound {
			fmt.Printf(domainScanLineFormat, v.Subdomain, v.Status, v.Ip)
		}

	}
	fmt.Printf("\nScanned %d domains in %.2fs.\n", len(scanResult.Subdomains), scanResult.Elapsed)
}

func (r JsonReporter) Report(result scan.Result) {
	marshal, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		return
	}

	fmt.Println(string(marshal))
}

func (r XmlReporter) Report(result scan.Result) {
	marshal, err := xml.MarshalIndent(result, "", " ")
	if err != nil {
		return
	}

	fmt.Println(string(marshal))
}

func NewReporter(flags utils.Flags) Reporter {
	switch flags.Format {
	case "json":
		return JsonReporter{flags: flags}
	case "xml":
		return XmlReporter{flags: flags}
	default:
		if flags.Type == scan.TypePort {
			return PortScanTextReporter{flags: flags}
		} else if flags.Type == scan.TypeDomain {
			return DomainScanTextReporter{flags: flags}
		}
	}

	return nil
}

func getColorByStatus(status string) string {
	switch status {
	case utils.PortOpen:
		return "\033[32m"
	default:
		return "\033[31m"
	}

}
