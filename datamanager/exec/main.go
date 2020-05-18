package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/honkytonktown/GoDistributed/datamanager"
	"github.com/honkytonktown/GoDistributed/dto"
	"github.com/honkytonktown/GoDistributed/qutils"
)

const url = "amqp://guest:guest@localhost:5672"

func main() {
	fmt.Println("Starting in main")
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		qutils.PersistReadingsQueue,
		"",
		false,
		true,
		false,
		false,
		nil)

	if err != nil {
		log.Fatalln("Failed to get access to messages")
	}
	fmt.Println("Consuming messages")
	for msg := range msgs {
		buf := bytes.NewReader(msg.Body)
		dec := gob.NewDecoder(buf)
		sd := &dto.SensorMessage{}
		dec.Decode(sd)
		fmt.Printf("The msg: %v\n", sd)

		err := datamanager.SaveReading(sd)
		if err != nil {
			log.Printf("Failed to save reading from sensor %v. Error: %s",
				sd.Name, err.Error())
		} else {
			msg.Ack(false)
		}
	}
}
