package datamanager

import (
	"fmt"

	"github.com/honkytonktown/GoDistributed/dto"
)

var sensors map[string]int

func SaveReading(reading *dto.SensorMessage) error {
	/* 	if sensors[reading.Name] == 0 {
	   		getSensors()
	   	}
	   	if sensors[reading.Name] == 0 {
	   		return errors.New("Unable to find sensor for name '" +
	   			reading.Name + "'")
	   	} */

	q := `
		INSERT INTO public."sensor_reading"
			(value, sensor_id, taken_on)
		VALUES
			($1, $2, $3)
		`

	_, err := db.Exec(q, reading.Value, 1, reading.Timestamp)

	return err
}

func getSensors() {
	sensors = make(map[string]int)
	q := `
		SELECT id, name
		FROM public."sensor"
		`

	rows, err := db.Query(q)
	fmt.Println("Before rows")
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		fmt.Println("Inside rows")
		var id int
		var name string
		rows.Scan(&id, &name)

		sensors[name] = id
		fmt.Printf("Id: %v, name: %v\n", id, name)
	}

}
