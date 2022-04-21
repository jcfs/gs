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

type PortScanTextReporter struct{}
type PortScanJsonReporter struct{}
type PortScanXmlReporter struct{}

//Report
func (r PortScanTextReporter) Report(result scan.Result) {
	scanResult := result.(scan.PortScanResult)

	fmt.Printf("%7s%7s%20s\n", "Port", "Status", "Service")
	fmt.Printf("%7s%7s%20s\n", "------", "------", "-------------------")
	for _, v := range scanResult.Ports {
		fmt.Printf("%7d\033[32m%7s\033[0m%20s\n", v.Port, v.Status, utils.GetPortDescription(v.Port))
	}
	fmt.Printf("\nCompleted in %.2fs.\n", scanResult.Elapsed)
}

func (r PortScanJsonReporter) Report(result scan.Result) {
	scanResult := result.(scan.PortScanResult)

	marshal, err := json.MarshalIndent(scanResult, "", " ")
	if err != nil {
		return
	}

	fmt.Println(string(marshal))
}

func (r PortScanXmlReporter) Report(result scan.Result) {
	scanResult := result.(scan.PortScanResult)

	marshal, err := xml.MarshalIndent(scanResult, "", " ")
	if err != nil {
		return
	}

	fmt.Println(string(marshal))
}

func NewReporter(flags utils.Flags) Reporter {
	if flags.Type == scan.TypePort {
		switch flags.Format {
		case "json":
			return PortScanJsonReporter{}
		case "xml":
			return PortScanXmlReporter{}
		default:
			return PortScanTextReporter{}
		}
	}

	return nil
}
