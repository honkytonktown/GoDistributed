package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/honkytonktown/GoDistributed/dto"
	"github.com/honkytonktown/GoDistributed/qutils"
	"github.com/streadway/amqp"
)

var url = "amqp://guest:guest@localhost:5672"

var name = flag.String("name", "sensor", "Name of the sensor")
var freq = flag.Uint("freq", 5, "Update frequency in cycles/sec")
var max = flag.Float64("max", 5.0, "Maximum value for generated readings")
var min = flag.Float64("min", 1.0, "Minimum value for generated readings")
var stepSize = flag.Float64("step", 0.1, "Maximum allowable change per measurement")

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var value = r.Float64()*(*max-*min) + *min
var nom = (*max-*min)/2 + *min

func main() {
	fmt.Println("Hello from sensors")
	flag.Parse()

	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch, false)
	publishQueueName(ch)
	discoveryQueue := qutils.GetQueue("", ch, true)
	ch.QueueBind(
		discoveryQueue.Name,
		"",
		qutils.SensorDiscoveryExchange,
		false,
		nil)

	go listenForDiscoveryRequests(discoveryQueue.Name, ch)
	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	signal := time.Tick(dur)

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	for range signal {
		calcValue()
		reading := dto.SensorMessage{
			Name:      *name,
			Value:     value,
			Timestamp: time.Now(),
		}
		buf.Reset()
		enc = gob.NewEncoder(buf)
		enc.Encode(reading)

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		ch.Publish(
			"",             //exchange name
			dataQueue.Name, //key string
			false,          //mandatory
			false,          //disable queue
			msg)            //sent message

		log.Printf("Reading sent. Value: %v\n", value)

	}
}
func listenForDiscoveryRequests(name string, ch *amqp.Channel) {
	msgs, _ := ch.Consume(
		name,
		"",
		true,
		false,
		false,
		false,
		nil)
	for range msgs {
		publishQueueName(ch)
	}
}
func publishQueueName(ch *amqp.Channel) {
	msg := amqp.Publishing{Body: []byte(*name)}
	ch.Publish(
		"amq.fanout", //exchange name
		"",           //key string
		false,        //mandatory
		false,        //disable queue
		msg)          //sent message

}
func calcValue() {
	var maxStep, minStep float64

	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	value += r.Float64()*(maxStep-minStep) + minStep
}
