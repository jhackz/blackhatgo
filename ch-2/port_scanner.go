package main

import (
	"fmt"
	"net"
	"sort"
)

// Worker func used for concurently scanning ports
// Worker accepts two channels.
func worker(ports, results chan int) {
	for p := range ports {
		addr := fmt.Sprintf("scanme.nmap.org:%d", p)

		// Dial connects to an RPC server at the specified network address.
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			// If the port doesn't respond send 0 back
			results <- 0
			continue
		}
		conn.Close()
		// Send the alive ports back accross the channel
		results <- p
	}
}

func main() {

	// make a ports channel intialized with a buffered size of 100
	// initialize and unbound results channel
	ports := make(chan int, 100)
	results := make(chan int)
	var open_ports []int

	// for the capacity of the ports channel
	// spin up the worker gorountine
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	// seperate go rountine for sending ports accross the channel for the worker
	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	// The result gathering loop receives on the results channel
	// We append the results to the open ports slice
	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			open_ports = append(open_ports, port)
		}
	}
	// Close our open channels
	close(ports)
	close(results)

	// Sort our results for display
	sort.Ints(open_ports)

	for _, port := range open_ports {
		fmt.Printf("[+] %d open\n", port)
	}

}
