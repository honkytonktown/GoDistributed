package dto

import (
	"encoding/gob"
	"time"
)

//SensorMesage is struct for the message that'll be sent to broker
type SensorMessage struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

func init() {
	gob.Register(SensorMessage{})
}
