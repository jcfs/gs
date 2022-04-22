package scan

import (
	"fmt"
	"gs/internal/utils"
	"net"
	"sync"
	"time"
)

type PortScanner struct{}

type PortStatus struct {
	Port   int    `json:"port,omitempty"`
	Status string `json:"status,omitempty"`
}

type PortScanResult struct {
	Elapsed float64      `json:"elapsed"`
	Ports   []PortStatus `json:"ports"`
}

func (scanner *PortScanner) Scan(flags utils.Flags, wg *sync.WaitGroup) Result {
	var (
		chwg       sync.WaitGroup
		portChunks = utils.Chunks(flags.Port, 250)
		statusChan = make(chan PortStatus)
		resultChan = make(chan PortScanResult)
		done       = make(chan int)
	)

	go output(statusChan, resultChan, done, wg)
	for _, chunk := range portChunks {
		for _, port := range chunk {
			scan(flags.Domain, port, statusChan, &chwg)
		}
		chwg.Wait()
	}

	done <- 1
	return <-resultChan
}

//scan async connect scan to a single port
func scan(host string, port int, result chan PortStatus, chwg *sync.WaitGroup) {
	chwg.Add(1)
	go func() {
		d := net.Dialer{Timeout: time.Duration(1000) * time.Millisecond}
		if _, err := d.Dial("tcp", fmt.Sprintf("%s:%d", host, port)); err == nil {
			result <- PortStatus{Port: port, Status: utils.PortOpen}
		} else {
			result <- PortStatus{Port: port, Status: utils.PortClosed}
		}
		chwg.Done()
	}()
}

func output(result chan PortStatus, resultChan chan PortScanResult, done chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	var (
		ports []PortStatus
		start = time.Now()
	)
	for {
		select {
		case r := <-result:
			ports = append(ports, r)
		case <-done:
			resultChan <- PortScanResult{Ports: ports, Elapsed: time.Since(start).Seconds()}
			wg.Done()
		}
	}
}
