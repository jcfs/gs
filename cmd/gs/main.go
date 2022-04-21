package gs

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

const banner = "GScanner 0.0.1-a (https://github.com/jcfs/gscanner)\n\n"

func main() {
	fmt.Print(banner)

	flags := Parse()
	if err := Validate(flags); err != nil {
		os.Exit(1)
	}

	scanner := NewScanner(flags)
	if scanner == nil {
		os.Exit(1)
	}

	start := time.Now()
	fmt.Printf("> Starting scanning [%s] at %v\n\n", flags.Type, time.Now())
	scanner.Scan(flags)
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("\n< Stopped scanning - %v\n", elapsed)
}
