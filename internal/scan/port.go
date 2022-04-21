package scan

import (
	"fmt"
	"gs/internal/utils"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type PortScanner struct{}

type Result struct {
	port   int
	status string
}

func (scanner *PortScanner) Scan(flags utils.Flags, wg *sync.WaitGroup) {
	var chwg sync.WaitGroup

	ch := utils.Chunks(flags.Port, 250)
	result := make(chan Result)
	done := make(chan int)

	go output(result, done, wg)
	for _, chunk := range ch {
		for _, port := range chunk {
			rawConnect(flags.Domain, port, result, &chwg)
		}
		chwg.Wait()
	}

	done <- 1
}

func rawConnect(host string, port int, result chan Result, chwg *sync.WaitGroup) {
	chwg.Add(1)
	go func() {
		d := net.Dialer{Timeout: time.Duration(150) * time.Millisecond}
		if _, err := d.Dial("tcp", host+":"+strconv.Itoa(port)); err != nil {
			if addErr, ok := err.(*net.AddrError); ok {
				if addErr.Timeout() {
					chwg.Done()
					return
				}
			} else {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			chwg.Done()
			return
		}

		result <- Result{port: port, status: utils.Open}
		chwg.Done()
	}()
}

func output(result chan Result, done chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	fmt.Printf("%7s%12s%20s\n", "Port", "Status", "Service")
	for {
		select {
		case r := <-result:
			fmt.Printf("%7d\033[32m%12s\033[0m%20s\n", r.port, r.status, utils.GetPortDescription(r.port))
		case <-done:
			wg.Done()
			return
		}
	}
}
