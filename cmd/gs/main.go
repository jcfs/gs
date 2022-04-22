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
	fmt.Printf("Starting scanning [%s: %s] at %v\n", flags.Type, flags.Domain, time.Now().Format("2006-02-01 15:04:05"))
	reporter := report.NewReporter(flags)
	reporter.Report(scanner.Scan(flags, &wg))
	wg.Wait()
}
