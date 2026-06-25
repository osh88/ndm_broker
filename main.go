// Решение тестового задания - брокер сообщений.
// go test -v ./...
// go run main.go -port 8080

package main

import (
	"flag"
	"log"
	"ndm_broker/broker"
	"ndm_broker/server"
)

var args struct {
	port              int
	initQueueCapacity int
}

func main() {
	flag.IntVar(&args.port, "port", 8080, "Broker API port")
	flag.IntVar(&args.initQueueCapacity, "initQueueCapacity", 50, "Init queue capacity")
	flag.Parse()

	switch {
	case args.port < 1 || 65535 < args.port:
		log.Fatalf("port < 1 or > 65535 (%d)", args.port)
	case args.initQueueCapacity < 1:
		log.Fatalf("initQueueCapacity < 1 (%d)", args.initQueueCapacity)
	}

	brk, err := broker.New(args.initQueueCapacity)
	check(err)

	srv := server.New(brk)
	check(srv.Run(args.port))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
