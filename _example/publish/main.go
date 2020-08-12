package main

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
	line "github.com/sminamot/nats-line-notify"
)

var (
	lineAccessToken string
)

func init() {
	lineAccessToken = os.Getenv("LINE_ACCESS_TOKEN")

	if lineAccessToken == "" {
		log.Fatalln("specify LINE_ACCESS_TOKEN")
	}
}

func main() {
	nc, err := nats.Connect("my-nats")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer ec.Close()

	// Publish the message
	//if err := ec.Publish("nats-line-notify", &line.Line{Message: "Hello, World!", ImageURL: "https://2.bp.blogspot.com/-xIRko5KAKaQ/WRaTiAqk1uI/AAAAAAABER0/QD8MxzLrmCwNNWOtiuNf54egnEwpQD7dACLcB/s400/vr_sweets_pokki_game.png", AccessToken: lineAccessToken}); err != nil {
	if err := ec.Publish("nats-line-notify", &line.Line{Message: "Hello, World!", AccessToken: lineAccessToken}); err != nil {
		log.Fatal(err)
	}
}
