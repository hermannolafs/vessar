package vessel

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

func DeclareChannelAndQueue(stopChannel chan os.Signal) (*amqp.Channel, amqp.Queue) {

	heartConnection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("%s: %s", "Fuck", err)
	}

	heartChannel, err := heartConnection.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Shit", err)
	}

	queHello, err := heartChannel.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	go WaitForConnectionAndChannelToClose(stopChannel, heartConnection, heartChannel)

	return heartChannel, queHello
}

func WaitForConnectionAndChannelToClose(
	stopChannel chan os.Signal,
	connection *amqp.Connection,
	channel *amqp.Channel,
) {
	<-stopChannel      // wait for SIGINT

	if err := channel.Close(); err != nil {
		log.Fatalf("%s: %s", "Shit", err)
	}
	if err := connection.Close(); err != nil {
		log.Fatalf("%s: %s", "Shit", err)
	}
}