package main

import (
	"log"
	"time"

	plugin_grpc "pilot-plugin/grpc"
)

const (
	servName = "Node manager"
)

func main() {
	log.Println("Start ", servName)

	plugin_grpc.StartServer()

	for {
		// TODO: Handle SIGTERM, Shutdown gracefully.
		time.Sleep(time.Second * 10)
	}
}
