package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/nats-io/nats.go"
	line "github.com/sminamot/nats-line-notify"
)

var (
	natsServer     string
	natsChannel    string
	natsQueueGroup string
)

func init() {
	natsServer = os.Getenv("NATS_SERVER")
	natsChannel = os.Getenv("NATS_CHANNEL")
	natsQueueGroup = os.Getenv("NATS_GROUP")

	switch "" {
	case natsServer, natsChannel, natsQueueGroup:
		log.Fatalln("specify environment variable")
	}

}

func main() {
	nc, err := nats.Connect(natsServer,
		nats.ErrorHandler(func(nc *nats.Conn, s *nats.Subscription, err error) {
			if s != nil {
				log.Printf("Async error in %q/%q: %v", s.Subject, s.Queue, err)
			} else {
				log.Printf("Async error outside subscription: %v", err)
			}
		}))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer ec.Close()

	// Subscribe
	// Decoding errors will be passed to the function supplied via
	// nats.ErrorHandler above, and the callback supplied here will
	// not be invoked.
	if _, err := ec.QueueSubscribe(natsChannel, natsQueueGroup, subscribeFunc); err != nil {
		log.Fatal(err)
	}

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	// Setup the interrupt handler to drain so we don't miss
	// requests when scaling down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	nc.Drain()
	log.Fatalf("Exiting")
}

func subscribeFunc(s *line.Line) {
	s.Notify()
}
