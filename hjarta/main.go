package main

import (
	"github.com/streadway/amqp"

	"hermannolafs/vessar/beinagrind"
	"hermannolafs/vessar/vessel"

	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	leikur := beinagrind.Leikur{
		Hnit: beinagrind.Hnit{PosX: 1, PosY: 2},
	}

	stopChannel := make(chan os.Signal)
	signal.Notify(stopChannel, os.Interrupt)

	heartChannel, queHello := vessel.DeclareChannelAndQueue(stopChannel)

	forever := make(chan bool)

	log.Printf("this is the leikur: %v", leikur)

	body, _ := leikur.ToBytes()

	log.Printf("this is the bytes: %v", body)

	go func() {
		for {
			err := heartChannel.Publish(
				"",     // exchange
				queHello.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        body,
				})
			failOnError(err, "Failed to publish a message")

			time.Sleep(time.Second * 2)
		}
	}()

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
