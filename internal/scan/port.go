package scan

import (
	"fmt"
	"gs/internal/utils"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PortScanner struct{}

var chwg sync.WaitGroup

type Result struct {
	port   int
	status string
}

func (scanner *PortScanner) Scan(flags utils.Flags, wg *sync.WaitGroup) {
	wg.Add(1)

	ch := utils.Chunks(flags.Port, 75)
	result := make(chan Result)
	done := make(chan int)

	go output(result, done, wg)
	for _, chunk := range ch {
		for _, port := range chunk {
			rawConnect(flags.Domain, port, result)
		}
		chwg.Wait()
	}

	done <- 1
	wg.Done()
}

func rawConnect(host string, port int, result chan Result) {
	chwg.Add(1)
	var timeoutTCP = time.Duration(100) * time.Millisecond
	go func() {
		d := net.Dialer{Timeout: timeoutTCP}
		_, err := d.Dial("tcp", host+":"+strconv.Itoa(port))
		if err != nil {
			if addErr, ok := err.(*net.AddrError); ok {
				if addErr.Timeout() {
					chwg.Done()
					return

				}
			} else if addErr, ok := err.(*net.OpError); ok {
				// handle lacked sufficient buffer space error
				if strings.TrimSpace(addErr.Err.Error()) == "bind: An operation on a socket could not be performed because "+
					"the system lacked sufficient buffer space or because a queue was full." {

					time.Sleep(timeoutTCP + (3000 * time.Millisecond))

					_, errAe := d.Dial("tcp", host+":"+strconv.Itoa(port))

					if errAe != nil {
						if addErr, ok := err.(*net.AddrError); ok {
							if addErr.Timeout() {
								chwg.Done()
								return

							}
						}
					}
				}

			} else {
				println(err.Error())
				os.Exit(1)

			}
			chwg.Done()
			return

		}

		result <- Result{port: port, status: "OPEN"}
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
