package main

import (
	"fmt"
	"gs/internal/common"
	"gs/internal/report"
	"gs/internal/scan"
	"gs/internal/utils"
	"os"
	"sync"
	"time"
)

const banner = "GScanner 0.0.1 (https://github.com/jcfs/gs)\n\n"

func main() {

	flags := utils.Parse(os.Args[1:])
	if err := common.Validate(flags); err != nil {
		os.Exit(1)
	}

	scanner := scan.NewScanner(flags)
	if scanner == nil {
		os.Exit(1)
	}

	var wg sync.WaitGroup
	fmt.Printf("Starting scanning [%s] at %v\n", flags.Type, time.Now().Format("2006-02-01 15:04:05"))
	result := scanner.Scan(flags, &wg)
	reporter := report.NewReporter(flags)
	reporter.Report(result)
	wg.Wait()
}
