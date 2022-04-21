package main

import (
	"fmt"
	"gs/internal/common"
	"gs/internal/scan"
	"gs/internal/utils"
	"os"
	"sync"
	"time"
)

const banner = "GScanner 0.0.1 (https://github.com/jcfs/gs)\n\n"

func main() {
	fmt.Print(banner)

	flags := utils.Parse(os.Args[1:])
	if err := common.Validate(flags); err != nil {
		os.Exit(1)
	}

	scanner := scan.NewScanner(flags)
	if scanner == nil {
		os.Exit(1)
	}

	var wg sync.WaitGroup
	start := time.Now()
	fmt.Printf("Starting scanning [%s] at %v\n", flags.Type, time.Now().Format("2006-02-01 15:04:05"))
	scanner.Scan(flags, &wg)
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("\nCompleted in %.2fs.\n", elapsed.Seconds())
}
