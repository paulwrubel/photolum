package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/paulwrubel/photolum/api"
)

func main() {
	fmt.Println("Starting Photolum...")

	fmt.Println("Starting API Server...")
	api.ListenAndServe()

	fmt.Println("Blocking until signalled to shutdown...")
	// make channel for interrupt signal
	c := make(chan os.Signal, 1)
	// tell os to send to chan when signal received
	signal.Notify(c, os.Interrupt)
	// wait for signal
	<-c

	fmt.Println("Shutting down...")
	os.Exit(0)
}
