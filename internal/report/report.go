package report

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gs/internal/scan"
	"gs/internal/utils"
	"io"
	"os"
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

type BaseReporter struct {
	flags utils.Flags
	w     io.Writer
}

type PortScanTextReporter struct {
	BaseReporter
}

type DomainScanTextReporter struct {
	BaseReporter
}

type JsonReporter struct {
	BaseReporter
}

type XmlReporter struct {
	BaseReporter
}

func (r PortScanTextReporter) Report(result scan.Result) {
	scanResult := result.(scan.PortScanResult)

	write(r.w, portScanHeaderFormat, "Port", "Status", "Service")
	write(r.w, portScanHeaderFormat, "------", "------", "-------------------")
	for _, v := range scanResult.Ports {
		color := utils.GetColorByStatus(v.Status)

		if r.flags.Verbose && v.Status != utils.PortOpen {
			write(r.w, portScanLineFormat, v.Port, color, v.Status, utils.GetPortDescription(v.Port))
		}

		if v.Status == utils.PortOpen {
			write(r.w, portScanLineFormat, v.Port, color, v.Status, utils.GetPortDescription(v.Port))
		}

	}
	write(r.w, "\nScanned %d ports in %.2fs.\n", len(scanResult.Ports), scanResult.Elapsed)
}

func (r DomainScanTextReporter) Report(result scan.Result) {
	scanResult := result.(scan.DomainScanResult)

	write(r.w, domainScanHeaderFormat, "Domain", "Status", "Service")
	write(r.w, domainScanHeaderFormat, "-------------------------", "----------", "----------------")
	for _, v := range scanResult.Subdomains {
		if r.flags.Verbose && v.Status != utils.DomainFound {
			write(r.w, domainScanLineFormat, v.Subdomain, v.Status, v.Ip)
		}

		if v.Status == utils.DomainFound {
			write(r.w, domainScanLineFormat, v.Subdomain, v.Status, v.Ip)
		}

	}
	write(r.w, "\nScanned %d domains in %.2fs.\n", len(scanResult.Subdomains), scanResult.Elapsed)
}

func (r JsonReporter) Report(result scan.Result) {
	marshal, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		return
	}

	write(r.w, string(marshal))
}

func (r XmlReporter) Report(result scan.Result) {
	marshal, err := xml.MarshalIndent(result, "", " ")
	if err != nil {
		return
	}

	write(r.w, string(marshal))
}

func NewReporter(flags utils.Flags) Reporter {
	var writer io.Writer = os.Stdout

	if flags.FileOutput != "" {
		if open, err := os.OpenFile(flags.FileOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); err == nil {
			writer = open
		} else {
			fmt.Println(err)
		}
	}

	baseReporter := BaseReporter{w: writer, flags: flags}

	switch flags.Format {
	case "json":
		return JsonReporter{baseReporter}
	case "xml":
		return XmlReporter{baseReporter}
	default:
		if flags.Type == scan.TypePort {
			return PortScanTextReporter{baseReporter}
		} else if flags.Type == scan.TypeDomain {
			return DomainScanTextReporter{baseReporter}
		}
	}

	return nil
}

func write(w io.Writer, format string, args ...any) {
	_, err := fmt.Fprintf(w, format, args...)
	if err != nil {
		return
	}
}
