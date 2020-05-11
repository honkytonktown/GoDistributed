package qutils

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const SensorListQueue = "SensorList"
const SensorDiscoveryExchange = "SensorDiscovery"
const PersistReadingsQueue = "PersistReading"

//GetChannel manages channel and connection for amqp
func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to establish a connection to the message broker")
	ch, err := conn.Channel()
	failOnError(err, "Failed to get channel for connection")

	return conn, ch

}

//GetQueue returns the amqp queue
func GetQueue(name string, ch *amqp.Channel, autoDelete bool) *amqp.Queue {
	q, err := ch.QueueDeclare(
		name,       //name of channel
		false,      //durable bool
		autoDelete, //autoDelete
		false,      //exclusivity
		false,      //noWait
		nil)
	failOnError(err, "Failed to declare queue")

	return &q

}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
