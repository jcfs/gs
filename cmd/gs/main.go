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

const banner = "GScanner 0.0.1-a (https://github.com/jcfs/gscanner)\n\n"

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
	fmt.Printf("> Starting scanning [%s] at %v\n\n", flags.Type, time.Now())
	scanner.Scan(flags, &wg)
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("\n< Stopped scanning [%v]\n", elapsed)
}
