package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func send(cancel context.CancelFunc, client TelnetClient) {
	defer cancel()
	if err := client.Send(); err != nil {
		fmt.Fprintf(os.Stderr, "unexpected sending err: %v", err)
		return
	}
	fmt.Fprintf(os.Stderr, "...EOF")
}

func receive(cancel context.CancelFunc, client TelnetClient) {
	defer cancel()
	if err := client.Receive(); err != nil {
		fmt.Fprintf(os.Stderr, "unexpected receiving error: %v", err)
		return
	}
	fmt.Fprintf(os.Stderr, "...connection was closed by peer")
}

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout for connection")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "undefined host and port")
		return
	}

	address := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	// if err := client.Connect(); err != nil {
	// 	fmt.Fprintf(os.Stderr, "unable to connect to server %s\n", address)
	// 	return
	// }

	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	go send(cancel, client)
	go receive(cancel, client)

	<-ctx.Done()
}
