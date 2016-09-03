package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/matiaskorhonen/oregon-go/oregonpi"
)

type thingState struct {
	SensorName       string  `json:"sensorName"`
	SensorType       int     `json:"sensorType"`
	SensorChannel    int     `json:"sensorChannel"`
	SensorLowBattery bool    `json:"sensorLowBattery"`
	Temperature      float32 `json:"temperature"`
	Humidity         float32 `json:"humidity"`
}

func main() {
	log.Println("Starting...")
	monitor, err := oregonpi.NewSensorMonitor(1, 0)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	readings := make(chan *oregonpi.SensorReading)
	monitor.ReadFromSensor(readings)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("Stopping...")
			monitor.Stop()
			close(readings)
		}
	}()

	for reading := range readings {
		log.Println(reading)
		updateThingShadow(reading)
	}
}

func updateThingShadow(reading *oregonpi.SensorReading) {
	sess, err := session.NewSession()
	if err != nil {
		log.Println("Failed to create AWS session: ", err)
		return
	}

	svc := iotdataplane.New(sess, aws.NewConfig().WithEndpoint(os.GetEnv("AWS_IOT_ENDPOINT")))

	reportedState := thingState{
		SensorName:       reading.Sensor.Name,
		SensorType:       reading.Sensor.Type,
		SensorChannel:    reading.Sensor.Channel,
		SensorLowBattery: reading.Sensor.LowBattery,
		Temperature:      reading.Temperature,
		Humidity:         reading.Humidity,
	}

	payload, err := json.Marshal(map[string]map[string]thingState{
		"state": map[string]thingState{
			"reported": reportedState,
		},
	})
	if err != nil {
		log.Println("Serialization error: ", err)
		return
	}

	log.Println("Updating Thing Shadow…")
	params := &iotdataplane.UpdateThingShadowInput{
		Payload:   payload,
		ThingName: aws.String("OutsideWeatherSensor"),
	}
	_, err := svc.UpdateThingShadow(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Println(err.Error())
		return
	}

	log.Println("Thing Shadow updated")
}
